package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware enregistre la méthode, le chemin et la durée de traitement.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Début requête %s %s", r.Method, r.URL.Path)
		log.Printf("En-têtes de la requête : %v", r.Header)
		next.ServeHTTP(w, r)
		log.Printf("En-têtes de la réponse : %v", w.Header())
		log.Printf("Fin requête %s %s en %v", r.Method, r.URL.Path, time.Since(start))
	})
}
