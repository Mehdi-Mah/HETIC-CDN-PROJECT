package loadbalancer

import "sync"

type LoadBalancer struct {
	servers []string
	index   int
	mu      sync.Mutex
}

// Instance singleton du load balancer
var lb *LoadBalancer

func Instance() *LoadBalancer {
	if lb == nil {
		// Liste de serveurs d'origine, à adapter selon l'infrastructure
		lb = &LoadBalancer{
			servers: []string{
				"http://origin-server-1:80",
				"http://origin-server-2:80",
			},
			index: 0,
		}
	}
	return lb
}

// NextServer sélectionne le serveur suivant en mode Round Robin.
func (lb *LoadBalancer) NextServer() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	server := lb.servers[lb.index]
	lb.index = (lb.index + 1) % len(lb.servers)
	return server
}
