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

	// Récupérer les données de la requête
	var req struct {
		Path string `json:"path"`
		Name string `json:"name"`
		Type string `json:"type"` // "file" ou "directory"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Vérification des paramètres
	if req.Path == "" || req.Name == "" || (req.Type != "file" && req.Type != "directory") {
		http.Error(w, "Chemin, nom ou type requis", http.StatusBadRequest)
		return
	}

	// Normalisation des chemins et noms de fichiers
	req.Path = strings.TrimSpace(req.Path)
	req.Path = filepath.Clean(req.Path)
	req.Name = strings.TrimSpace(req.Name)

	// Construction du chemin sécurisé
	basePath := fmt.Sprintf("uploads/%sUploads", username)

	// Pour un fichier, on ajoute le nom du fichier à la fin du chemin
	var fullPath string
	if req.Type == "file" {
		fullPath = filepath.Join(basePath, req.Path, req.Name)
	} else {
		fullPath = filepath.Join(basePath, req.Path)
	}
	fullPath = filepath.Clean(fullPath)

	// Logs pour le débogage
	fmt.Println("Base Path:", basePath)
	fmt.Println("Chemin Reçu:", req.Path)
	fmt.Println("Nom Reçu:", req.Name)
	fmt.Println("Chemin Final:", fullPath)

	// Vérifier que l'utilisateur ne supprime que ses propres fichiers/dossiers
	if !strings.HasPrefix(fullPath, basePath) {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	// Vérifier si le fichier/dossier existe
	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		fmt.Println("Fichier ou dossier introuvable:", fullPath)
		http.Error(w, "Fichier ou dossier introuvable", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("Erreur d'accès au fichier:", err)
		http.Error(w, "Erreur lors de l'accès au fichier", http.StatusInternalServerError)
		return
	}

	// Suppression du fichier ou dossier
	if req.Type == "file" {
		// Vérifier que ce n'est pas un dossier avant suppression
		if fileInfo.IsDir() {
			fmt.Println("Erreur: tentative de suppression d'un dossier en tant que fichier.")
			http.Error(w, "Impossible de supprimer un dossier en tant que fichier", http.StatusBadRequest)
			return
		}

		// Supprimer uniquement le fichier
		err = os.Remove(fullPath)
		if err != nil {
			fmt.Println("Erreur suppression fichier:", err)
			http.Error(w, "Erreur lors de la suppression du fichier", http.StatusInternalServerError)
			return
		}

		fmt.Println("Fichier supprimé avec succès:", fullPath)

	} else if req.Type == "directory" {
		// Vérifier que c'est bien un dossier avant suppression
		if !fileInfo.IsDir() {
			fmt.Println("Erreur: tentative de suppression d'un fichier en tant que dossier.")
			http.Error(w, "Impossible de supprimer un fichier en tant que dossier", http.StatusBadRequest)
			return
		}

		// Supprimer le dossier et tout son contenu
		err = os.RemoveAll(fullPath)
		if err != nil {
			fmt.Println("Erreur suppression dossier:", err)
			http.Error(w, "Erreur lors de la suppression du dossier", http.StatusInternalServerError)
			return
		}

		fmt.Println("Dossier et son contenu supprimés avec succès:", fullPath)
	}

	// Réponse JSON
	json.NewEncoder(w).Encode(map[string]string{"message": "Suppression réussie"})
}