package proxy

import (
	"HETIC-CDN-PROJECT/pkg/loadbalancer"
	"io"
	"log"
	"net/http"
)

// NewProxyHandler retourne un handler HTTP pour le proxy.
func NewProxyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Sélection du serveur d'origine via le load balancer
		target := loadbalancer.Instance().NextServer()
		proxyURL := target + r.RequestURI

		log.Printf("Redirection de la requête %s vers %s", r.RequestURI, proxyURL)

		// Création de la nouvelle requête à rediriger
		req, err := http.NewRequest(r.Method, proxyURL, r.Body)
		if err != nil {
			http.Error(w, "Erreur lors de la création de la requête", http.StatusInternalServerError)
			return
		}

		// Copie des en-têtes de la requête d'origine
		req.Header = r.Header

		// Envoi de la requête vers le serveur d'origine
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Erreur de communication avec le serveur d'origine", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copie des en-têtes et du code de réponse
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})
}
