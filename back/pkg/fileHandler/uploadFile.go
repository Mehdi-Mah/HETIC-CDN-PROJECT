package fileHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UploadHandler gère l'upload des fichiers
func UploadHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("🔹 Requête reçue pour upload")

    // Vérification du token
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        fmt.Println("Token manquant !")
        http.Error(w, "Token manquant", http.StatusUnauthorized)
        return
    }

    username, err := extractUsernameFromToken(authHeader)
    if err != nil {
        fmt.Println("Token invalide !")
        http.Error(w, "Token invalide", http.StatusUnauthorized)
        return
    }
    fmt.Println("Utilisateur authentifié :", username)

    // Vérification du fichier
    file, header, err := r.FormFile("file")
    if err != nil {
        fmt.Println("Erreur récupération fichier :", err)
        http.Error(w, "Fichier non reçu", http.StatusBadRequest)
        return
    }
    defer file.Close()
    fmt.Println("Fichier reçu :", header.Filename)

    // Récupération du chemin de destination
    uploadPath := r.FormValue("path")
    if uploadPath == "" {
        uploadPath = "/"
    }
    fmt.Println("📂 Chemin d'upload :", uploadPath)

    // Création du chemin complet
    basePath := fmt.Sprintf("uploads/%sUploads", username)
    fullPath := filepath.Join(basePath, uploadPath, header.Filename)
    fullPath = filepath.Clean(fullPath)

    // Vérification de la sécurité du chemin
    if !strings.HasPrefix(fullPath, basePath) {
        fmt.Println("Tentative d'accès non autorisée :", fullPath)
        http.Error(w, "Chemin invalide", http.StatusForbidden)
        return
    }

    // Création du dossier si inexistant
    if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
        fmt.Println("Erreur lors de la création du dossier :", err)
        http.Error(w, "Erreur de création du dossier", http.StatusInternalServerError)
        return
    }

    // Création du fichier sur le serveur
    dst, err := os.Create(fullPath)
    if err != nil {
        fmt.Println("Erreur de création du fichier :", err)
        http.Error(w, "Erreur lors de la création du fichier", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    // Copie du contenu du fichier
    if _, err := io.Copy(dst, file); err != nil {
        fmt.Println("Erreur de copie du fichier :", err)
        http.Error(w, "Erreur lors de l'écriture du fichier", http.StatusInternalServerError)
        return
    }

    fmt.Println("Fichier uploadé avec succès :", fullPath)

    // Réponse au frontend
    json.NewEncoder(w).Encode(map[string]string{"message": "Fichier uploadé avec succès", "path": fullPath})
}

