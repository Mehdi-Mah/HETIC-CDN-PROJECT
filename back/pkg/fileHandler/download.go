package fileHandler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DownloadHandler gère le téléchargement d'un fichier
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Vérifier le token pour l'authentification
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token manquant", http.StatusUnauthorized)
		return
	}

	username, err := extractUsernameFromToken(authHeader)
	if err != nil {
		http.Error(w, "Token invalide", http.StatusUnauthorized)
		return
	}

	// Récupérer les paramètres `path` et `name` du fichier à télécharger
	filePath := r.URL.Query().Get("path")
	fileName := r.URL.Query().Get("name")

	if filePath == "" || fileName == "" {
		http.Error(w, "Chemin et nom du fichier requis", http.StatusBadRequest)
		return
	}

	// Normalisation des chemins et noms de fichiers
	filePath = strings.TrimSpace(filePath)
	filePath = filepath.Clean(filePath)
	fileName = strings.TrimSpace(fileName)

	// Construire le chemin absolu du fichier
	basePath := fmt.Sprintf("uploads/%sUploads", username)
	fullPath := filepath.Join(basePath, filePath, fileName)
	fullPath = filepath.Clean(fullPath)

	// Logs pour le débogage
	fmt.Println("Base Path:", basePath)
	fmt.Println("Chemin Reçu:", filePath)
	fmt.Println("Nom du Fichier Reçu:", fileName)
	fmt.Println("Chemin Final:", fullPath)

	// Vérifier que le fichier appartient bien à l'utilisateur
	if !strings.HasPrefix(fullPath, basePath) {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	// Vérifier que le fichier existe
	file, err := os.Open(fullPath)
	if err != nil {
		fmt.Println("Erreur: fichier introuvable", fullPath)
		http.Error(w, "Fichier introuvable", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Définir les en-têtes pour forcer le téléchargement
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Envoyer le fichier en réponse
	http.ServeFile(w, r, fullPath)
}
