package fileHandler

import (
	"net/http"
	"os"
	"path/filepath"
)

// DownloadHandler gère le téléchargement des fichiers
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Paramètre 'file' manquant", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadPath, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}
