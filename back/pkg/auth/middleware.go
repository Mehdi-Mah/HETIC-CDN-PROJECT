package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware vérifie le token JWT et ajoute l'ID utilisateur dans le contexte.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "En-tête d'authentification manquant", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "En-tête d'authentification invalide", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Vérifie que la méthode de signature est HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Token invalide", http.StatusUnauthorized)
			return
		}
		// Ajoute l'ID utilisateur dans le contexte, si besoin pour les handlers ultérieurs
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userId, ok := claims["userId"].(string); ok {
				ctx := context.WithValue(r.Context(), "userId", userId)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}
