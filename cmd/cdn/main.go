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
	// Liste des serveurs d'origine pour le reverse proxy
	targets := []string{
		"http://cdn-project-api:8080", // votre API principale dockerrise
		"http://localhost:8080",       // api local
	}

	// Création de l'instance du reverse proxy avec failover
	reverseProxy := proxy.NewFailoverReverseProxy(targets)

	// Récupérer l'URI MongoDB depuis l'environnement
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Valeur par défaut pour le développement hors conteneur
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:password@localhost:27017/cdnproject?authSource=admin"))
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

	// Création du multiplexer principal
	mux := http.NewServeMux()

	// Endpoints d'authentification
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	// Intégration du reverse proxy :
	// Les requêtes commençant par /proxy seront redirigées via le reverse proxy.
	mux.Handle("/proxy/", http.StripPrefix("/proxy", reverseProxy))

	// Route basique pour vérifier la disponibilité du service
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Routes pour l'upload et le téléchargement de fichiers
	mux.HandleFunc("/upload", fileHandler.UploadHandler)
	mux.HandleFunc("/download", fileHandler.DownloadHandler)

	// Application des middlewares CORS et de logging
	muxWithMiddleware := middleware.LoggingMiddleware(corsMiddleware(mux))

	// Configuration du serveur avec timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      muxWithMiddleware,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Démarrage du serveur en HTTPS si configuré, sinon en HTTP
	if security.UseTLS() {
		log.Println("Serveur démarré en HTTPS sur le port 8080")
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		log.Println("Serveur démarré en HTTP sur le port 8080")
		log.Fatal(server.ListenAndServe())
	}
}
