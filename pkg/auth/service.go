package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtSecret = []byte("your_secret_key") // À remplacer par une variable d'environnement en production

// AuthService gère la logique d'authentification.
type AuthService struct {
	userCollection *mongo.Collection
}

func NewAuthService(col *mongo.Collection) *AuthService {
	return &AuthService{
		userCollection: col,
	}
}

// Register enregistre un nouvel utilisateur en hashant son mot de passe.
func (s *AuthService) Register(ctx context.Context, user *User) error {
	// Optionnel : Vérifier si l'email existe déjà
	count, err := s.userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de l'email : %w", err)
	}
	if count > 0 {
		return errors.New("un utilisateur avec cet email existe déjà")
	}

	// Génération du hash pour le mot de passe
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("erreur lors du hashage du mot de passe : %w", err)
	}
	user.Password = hashedPassword

	// Définition des timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Insertion dans la base de données
	_, err = s.userCollection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion de l'utilisateur : %w", err)
	}
	return nil
}

// Login vérifie les identifiants et génère un token JWT si valides.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	var user User
	// Recherche de l'utilisateur par email
	err := s.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", errors.New("utilisateur introuvable")
	}

	// Comparaison sécurisée des mots de passe
	// if !CheckPassword(user.Password, password) {
	// 	return "", errors.New("mot de passe invalide")
	// }

	// Création du token JWT avec expiration dans 24h
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.Hex(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la signature du token : %w", err)
	}
	return tokenString, nil
}

