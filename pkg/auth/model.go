package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User représente un utilisateur dans la base de données.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"` // Le mot de passe ne sera pas renvoyé dans les réponses JSON
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
