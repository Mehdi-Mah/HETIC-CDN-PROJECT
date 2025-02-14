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

	// Récupérer le paramètre `path` du fichier à télécharger
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Chemin du fichier requis", http.StatusBadRequest)
		return
	}

	// Construire le chemin absolu du fichier
	basePath := fmt.Sprintf("uploads/%sUploads", username)
	fullPath := filepath.Join(basePath, filePath)
	fullPath = filepath.Clean(fullPath)

	// Vérifier que le fichier appartient bien à l'utilisateur
	if !strings.HasPrefix(fullPath, basePath) {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	// Vérifier que le fichier existe
	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "Fichier introuvable", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Définir les en-têtes pour forcer le téléchargement
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(fullPath))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Envoyer le fichier en réponse
	http.ServeFile(w, r, fullPath)
}
