package fileHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// CreateFolderHandler gère la création d'un dossier avec une profondeur limitée
func CreateFolderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Vérification du token
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

	// Décoder la requête JSON
	var req struct {
		FolderName string `json:"folderName"`
		Path       string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	if req.FolderName == "" {
		http.Error(w, "Nom du dossier requis", http.StatusBadRequest)
		return
	}

	// Construire le chemin absolu
	basePath := fmt.Sprintf("uploads/%sUploads", username)
	fullPath := filepath.Join(basePath, req.Path, req.FolderName)
	fullPath = filepath.Clean(fullPath)

	// Vérifier la profondeur du chemin (limite de 10 niveaux)
	relPath := strings.TrimPrefix(fullPath, basePath)
	relPath = strings.TrimPrefix(relPath, "/")
	depth := len(strings.Split(relPath, "/"))

	if depth > 10 {
		http.Error(w, "Profondeur maximale atteinte (10 niveaux)", http.StatusForbidden)
		return
	}

	// Vérifier que l'utilisateur ne tente pas de sortir de son dossier
	if !strings.HasPrefix(fullPath, basePath) {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	// Créer le dossier si la profondeur est valide
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		http.Error(w, "Erreur lors de la création du dossier", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Dossier créé avec succès"})
}
