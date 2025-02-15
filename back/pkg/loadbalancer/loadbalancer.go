package loadbalancer

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Algorithm définit le type d'algorithme de répartition.
type Algorithm int

const (
	RoundRobin Algorithm = iota
	WeightedRoundRobin
	LeastConnections
)

// BackendServer représente un serveur backend.
type BackendServer struct {
	URL               *url.URL
	Weight            int        // Pour Weighted Round Robin
	ActiveConnections int64      // Pour Least Connections
	FailCount         int        // Nombre de vérifications de santé échouées
	Online            bool       // Statut du serveur
	Mutex             sync.Mutex // Pour synchroniser les mises à jour
}

// LoadBalancer gère la liste des serveurs et la répartition des requêtes.
type LoadBalancer struct {
	Servers       []*BackendServer
	algo          Algorithm
	rrIndex       uint64        // Compteur pour Round Robin
	timeout       time.Duration // Timeout pour les requêtes backend
	checkInterval time.Duration // Intervalle de health check
}

// NewLoadBalancer crée une instance de LoadBalancer avec la liste des cibles (URL en chaîne),
// l'algorithme choisi, un timeout et un intervalle de vérification.
func NewLoadBalancer(servers []string, algo Algorithm, timeout, checkInterval time.Duration) *LoadBalancer {
	lb := &LoadBalancer{
		algo:          algo,
		timeout:       timeout,
		checkInterval: checkInterval,
	}
	for _, s := range servers {
		u, err := url.Parse(s)
		if err != nil {
			log.Fatalf("Erreur lors de l'analyse de l'URL %s: %v", s, err)
		}
		lb.Servers = append(lb.Servers, &BackendServer{
			URL:    u,
			Weight: 1,    // Poids par défaut (modifiable dynamiquement si besoin)
			Online: true, // Par défaut, le serveur est considéré comme en ligne
		})
	}
	// Lancer les health checks en arrière-plan
	go lb.healthCheckLoop()
	return lb
}

// getNextServer sélectionne un serveur en fonction de l'algorithme choisi parmi ceux en ligne.
func (lb *LoadBalancer) getNextServer() (*BackendServer, error) {
	// Filtrer les serveurs en ligne
	var onlineServers []*BackendServer
	for _, s := range lb.Servers {
		if s.Online {
			onlineServers = append(onlineServers, s)
		}
	}
	if len(onlineServers) == 0 {
		return nil, errors.New("aucun serveur en ligne")
	}

	switch lb.algo {
	case RoundRobin:
		idx := int(atomic.AddUint64(&lb.rrIndex, 1) % uint64(len(onlineServers)))
		return onlineServers[idx], nil
	case WeightedRoundRobin:
		totalWeight := 0
		for _, s := range onlineServers {
			totalWeight += s.Weight
		}
		if totalWeight == 0 {
			return onlineServers[0], nil
		}
		rr := int(atomic.AddUint64(&lb.rrIndex, 1) % uint64(totalWeight))
		for _, s := range onlineServers {
			rr -= s.Weight
			if rr < 0 {
				return s, nil
			}
		}
		return onlineServers[0], nil
	case LeastConnections:
		var selected *BackendServer
		minConns := int64(^uint64(0) >> 1) // valeur max d'int64
		for _, s := range onlineServers {
			conns := atomic.LoadInt64(&s.ActiveConnections)
			if conns < minConns {
				minConns = conns
				selected = s
			}
		}
		if selected == nil {
			selected = onlineServers[0]
		}
		return selected, nil
	default:
		return nil, errors.New("algorithme non supporté")
	}
}

// ServeHTTP redirige la requête entrante vers un serveur backend choisi.
// En cas d'erreur (ex : serveur offline), il tente les autres serveurs (failover).
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	attempts := len(lb.Servers)
	fmt.Printf("attempts: %d\n", attempts)
	for i := 0; i < attempts; i++ {
		server, err := lb.getNextServer()
		if err != nil {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		// Incrémenter le compteur de connexions actives pour ce serveur
		atomic.AddInt64(&server.ActiveConnections, 1)

		// Créer un reverse proxy vers le serveur sélectionné
		proxy := httputil.NewSingleHostReverseProxy(server.URL)
		proxy.Transport = &http.Transport{
			ResponseHeaderTimeout: lb.timeout,
		}

		// Personnaliser le Director pour ajouter des en-têtes
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Header.Set("X-Forwarded-For", r.RemoteAddr)
		}

		// Variable de contrôle d'erreur pour cette tentative
		var failed bool
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Printf("Erreur avec le serveur %s: %v", server.URL, err)
			failed = true
			// Augmenter le compteur d'échecs et marquer le serveur offline si nécessaire
			server.Mutex.Lock()
			server.FailCount++
			if server.FailCount >= 3 {
				server.Online = false
				log.Printf("Serveur %s marqué hors service", server.URL)
			}
			server.Mutex.Unlock()
		}

		proxy.ServeHTTP(w, r)

		// Décrémenter le compteur de connexions actives
		atomic.AddInt64(&server.ActiveConnections, -1)

		// Si aucune erreur n'a été détectée, la requête a été traitée avec succès
		if !failed {
			return
		}
	}
	// Si toutes les tentatives échouent, renvoyer une erreur 503
	http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
}

// healthCheckLoop lance un health check sur chaque serveur à intervalle régulier.
func (lb *LoadBalancer) healthCheckLoop() {
	ticker := time.NewTicker(lb.checkInterval)
	for range ticker.C {
		for _, server := range lb.Servers {
			go lb.checkServerHealth(server)
		}
	}
}

// checkServerHealth vérifie la disponibilité d'un serveur (par exemple en appelant /health).
func (lb *LoadBalancer) checkServerHealth(server *BackendServer) {
	client := &http.Client{Timeout: lb.timeout}
	healthURL := server.URL.String() + "/health"
	resp, err := client.Get(healthURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		server.Mutex.Lock()
		server.FailCount++
		if server.FailCount >= 3 {
			server.Online = false
			log.Printf("Health check échoué pour %s, serveur hors service", server.URL)
			// Vous pouvez ajouter ici une notification (Slack, email, etc.)
		}
		server.Mutex.Unlock()
		if resp != nil {
			resp.Body.Close()
		}
		return
	}
	// Si le serveur répond, réinitialiser les compteurs et marquer comme online
	server.Mutex.Lock()
	if !server.Online {
		log.Printf("Serveur %s de nouveau en ligne", server.URL)
	}
	server.FailCount = 0
	server.Online = true
	server.Mutex.Unlock()
	resp.Body.Close()
}