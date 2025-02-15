package fileHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// FileInfo représente un fichier ou un dossier dans l'arborescence
type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"` // "file" ou "folder"
	Path string `json:"path"` // Chemin relatif
}

// ListFilesHandler gère la récupération de l'arborescence des fichiers d'un utilisateur
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Récupérer le token depuis les headers
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token manquant", http.StatusUnauthorized)
		return
	}

	// Extraire le username du token JWT
	username, err := extractUsernameFromToken(authHeader)
	if err != nil {
		http.Error(w, "Token invalide", http.StatusUnauthorized)
		return
	}

	// Définir le chemin du dossier de l'utilisateur
	userDir := fmt.Sprintf("uploads/%sUploads", username)

	// Vérifier si le dossier existe
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		json.NewEncoder(w).Encode([]FileInfo{}) // Retourne une liste vide si pas de fichiers
		return
	}

	// Récupérer l'arborescence des fichiers
	files := []FileInfo{}
	err = filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Chemin relatif par rapport au dossier utilisateur
		relPath := strings.TrimPrefix(path, userDir)
		relPath = strings.TrimPrefix(relPath, "/") // Enlever le / au début si nécessaire

		if relPath == "" {
			return nil // On ignore le dossier root de l'utilisateur
		}

		// Vérifier si c'est un fichier ou un dossier
		fileType := "file"
		cleanPath := relPath // Par défaut, garder le chemin tel quel

		if info.IsDir() {
			fileType = "folder"
		} else {
			// ✅ Ne garder que le dossier parent pour un fichier
			cleanPath = filepath.Dir(relPath)
			if cleanPath == "." { // Si le fichier est dans le dossier root, path devient ""
				cleanPath = ""
			}
		}

		files = append(files, FileInfo{
			Name: info.Name(),
			Type: fileType,
			Path: cleanPath,
		})

		return nil
	})

	if err != nil {
		http.Error(w, "Erreur lors de la récupération des fichiers", http.StatusInternalServerError)
		return
	}

	// Retourner les fichiers en JSON
	json.NewEncoder(w).Encode(files)
}

// extractUsernameFromToken extrait le username du token JWT
func extractUsernameFromToken(authHeader string) (string, error) {
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	fmt.Println("Token reçu:", tokenString) // 🔍 Log pour voir le token

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Println("Erreur de parsing du token:", err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("Claims du token:", claims) // 🔍 Voir les claims
		if username, exists := claims["username"].(string); exists {
			fmt.Println("Username extrait du token:", username)
			return username, nil
		}
	}
	fmt.Println("Erreur: username non trouvé dans le token")
	return "", fmt.Errorf("username non trouvé dans le token")
}