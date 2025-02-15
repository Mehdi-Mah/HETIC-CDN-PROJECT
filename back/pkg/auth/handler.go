package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// AuthHandler gère les requêtes HTTP liées à l'authentification.
type AuthHandler struct {
	Service *AuthService
}

func NewAuthHandler(col *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		Service: NewAuthService(col),
	}
}

// Register gère l'inscription d'un nouvel utilisateur.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Entrée invalide", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.Service.Register(ctx, &user); err != nil {
		http.Error(w, "Échec de l'inscription : "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Utilisateur inscrit", "username": user.Username})
}

// Login gère la connexion et renvoie un token JWT en cas de succès.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Entrée invalide", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	token, err := h.Service.Login(ctx, credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Échec de la connexion : "+err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
