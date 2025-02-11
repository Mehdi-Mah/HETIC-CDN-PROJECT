package main

import (
	"HETIC-CDN-PROJECT/pkg/middleware"
	"HETIC-CDN-PROJECT/pkg/proxy"
	"HETIC-CDN-PROJECT/pkg/security"
	"log"
	"net/http"
	"time"
)

func main() {
	/*
		Crée un multiplexer qui va gérer les différentes routes de l’application.
		 C’est ici on associe les URL à des fonctions spécifiques.*/
	mux := http.NewServeMux()
	muxWithMiddleware := middleware.LoggingMiddleware(mux)

	/* Route du proxy pour rediriger les requêtes, vers les serveurs d’origine.
	   Le package proxy va choisir le serveur via le load balancer*/
	mux.Handle("/", proxy.NewProxyHandler())

	// Ajout d'une route basique pour vérifier la disponibilité du service
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Configuration du serveur avec timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      muxWithMiddleware, // Utilisation du multiplexer avec middleware
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second, // Timeout pour les connexions inactives
	}

	// Démarrage du serveur en mode HTTPS si configuré, sinon en HTTP
	if security.UseTLS() {
		log.Println("Serveur démarré en HTTPS sur le port 8080")
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		log.Println("Serveur démarré en HTTP sur le port 8080")
		log.Fatal(server.ListenAndServe())
	}
}
