package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err = s.userCollection.InsertOne(ctx, user)
	return err
}

// Login vérifie les identifiants et génère un token JWT si valides.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	var user User
	err := s.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", errors.New("utilisateur introuvable")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("identifiants invalides")
	}

	// Création du token JWT avec une expiration de 24h
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.Hex(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
