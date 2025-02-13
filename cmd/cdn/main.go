package main

import (
	"HETIC-CDN-PROJECT/handler/fileHandler"
	"HETIC-CDN-PROJECT/pkg/auth"
	"HETIC-CDN-PROJECT/pkg/middleware"
	"HETIC-CDN-PROJECT/pkg/proxy"
	"HETIC-CDN-PROJECT/pkg/security"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Récupérer l'URI MongoDB depuis l'environnement
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Valeur par défaut pour le développement hors conteneur
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("cdnproject")
	userCollection := db.Collection("users")

	// Initialisation de l'handler d'authentification
	authHandler := auth.NewAuthHandler(userCollection)

	// Middleware pour gérer les en-têtes CORS
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	mux := http.NewServeMux()
	muxWithMiddleware := middleware.LoggingMiddleware(corsMiddleware(mux))

	// Endpoints d'authentification
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	// Route du proxy pour rediriger les requêtes vers les serveurs d’origine
	mux.Handle("/", proxy.NewProxyHandler())

	// Ajout d'une route basique pour vérifier la disponibilité du service
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Route pour l'upload de fichiers
	mux.HandleFunc("/upload", fileHandler.UploadHandler)
	// Route pour le téléchargement de fichiers
	mux.HandleFunc("/download", fileHandler.DownloadHandler)

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
