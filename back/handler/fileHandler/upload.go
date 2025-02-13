package fileHandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// UploadHandler gère l'upload des fichiers
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Analyse le formulaire multipart (limite fixée à 10 MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération du fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Vérifier si le fichier est vide en utilisant fileHeader.Size
	if fileHeader.Size == 0 {
		http.Error(w, "Le fichier est vide, upload refusé", http.StatusBadRequest)
		return
	}

	// S'assurer que le répertoire de destination existe
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		err = os.MkdirAll(uploadPath, os.ModePerm)
		if err != nil {
			http.Error(w, "Erreur lors de la création du répertoire", http.StatusInternalServerError)
			return
		}
	}

	dst, err := os.Create(filepath.Join(uploadPath, fileHeader.Filename))
	if err != nil {
		http.Error(w, "Erreur lors de la création du fichier", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Erreur lors de la copie du fichier", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Fichier uploadé avec succès : %s", fileHeader.Filename)
}
