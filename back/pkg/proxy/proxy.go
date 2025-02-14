package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

// FailoverReverseProxy gère la redirection vers plusieurs serveurs backend avec failover.
type FailoverReverseProxy struct {
	targets   []*url.URL      // Liste des URL des serveurs backend
	current   uint64          // Compteur atomique pour la sélection round robin
	transport *http.Transport // Transport personnalisé avec des timeouts
}

// NewFailoverReverseProxy construit le proxy en analysant une liste de cibles.
func NewFailoverReverseProxy(targets []string) *FailoverReverseProxy {
	var urls []*url.URL
	for _, target := range targets {
		parsedURL, err := url.Parse(target)
		if err != nil {
			log.Fatalf("Erreur lors de l'analyse de l'URL %s: %v", target, err)
		}
		urls = append(urls, parsedURL)
	}
	return &FailoverReverseProxy{
		targets: urls,
		transport: &http.Transport{
			ResponseHeaderTimeout: 5 * time.Second,
			IdleConnTimeout:       30 * time.Second,
		},
	}
}

// ServeHTTP redirige la requête entrante vers l'un des serveurs backend.
// En cas d'erreur sur une cible, il tente la suivante (failover) jusqu'à épuisement des cibles.
func (p *FailoverReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	numTargets := len(p.targets)

	// On essaie chaque cible une fois
	for i := 0; i < numTargets; i++ {
		// Sélection round robin : incrément atomique et modulo nombre de cibles
		idx := int(atomic.AddUint64(&p.current, 1) % uint64(numTargets))
		target := p.targets[idx]

		// Création d'un reverse proxy pour la cible sélectionnée
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = p.transport

		// Sauvegarde du Director original pour configurer l'URL de la cible
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			// Configure la requête pour pointer vers la cible
			originalDirector(req)
			// Ajoute un header X-Forwarded-For pour transmettre l'IP du client
			req.Header.Set("X-Forwarded-For", r.RemoteAddr)
		}

		// Variable de contrôle pour détecter une erreur sur cette cible
		failed := false

		// ErrorHandler qui loggue l'erreur et indique l'échec
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Printf("Erreur sur la cible %s: %v", target.String(), err)
			failed = true
		}

		// Log de la cible sélectionnée
		log.Printf("Redirection de la requête %s vers %s", r.RequestURI, target.String())

		// Utilisation du proxy pour relayer la requête
		proxy.ServeHTTP(w, r)

		// Si la requête a été traitée sans échec, on arrête ici
		if !failed {
			return
		}
	}

	// Si toutes les cibles échouent, renvoie une erreur 503 Service Unavailable
	http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
}
