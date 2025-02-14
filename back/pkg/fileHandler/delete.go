package fileHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DeleteHandler gère la suppression d'un fichier ou d'un dossier
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

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

	// Récupérer le chemin du fichier/dossier à supprimer depuis le corps de la requête
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(w, "Chemin du fichier/dossier requis", http.StatusBadRequest)
		return
	}

	// Construire le chemin absolu du fichier/dossier
	basePath := fmt.Sprintf("uploads/%sUploads", username)
	fullPath := filepath.Join(basePath, req.Path)
	fullPath = filepath.Clean(fullPath)

	// Vérifier que l'utilisateur ne supprime que ses propres fichiers/dossiers
	if !strings.HasPrefix(fullPath, basePath) {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	// Vérifier si le fichier/dossier existe
	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		http.Error(w, "Fichier ou dossier introuvable", http.StatusNotFound)
		return
	}

	// Supprimer le fichier ou le dossier
	if fileInfo.IsDir() {
		err = os.RemoveAll(fullPath) // Supprime le dossier et son contenu
	} else {
		err = os.Remove(fullPath) // Supprime seulement le fichier
	}

	if err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}

	// Réponse JSON
	json.NewEncoder(w).Encode(map[string]string{"message": "Suppression réussie"})
}
