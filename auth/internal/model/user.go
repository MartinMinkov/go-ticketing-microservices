package model

import (
	"errors"
	"time"

	"auth.mminkov.net/internal/database"
	"auth.mminkov.net/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Password   *string            `json:"-" validate:"required,min=6"`
	Email      *string            `json:"email" validate:"email,required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}

func NewUser(email string, password string) *User {
	return &User{
		ID:         primitive.NewObjectID(),
		Email:      &email,
		Password:   &password,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}
}

func FindByEmail(db *database.Database, email string) (*User, error) {
	var user User
	err := db.UserCollection.FindOne(db.Ctx, bson.D{{Key: "email", Value: email}}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, errors.New("failed to find user: " + err.Error())
	}
	return &user, nil
}

func (u *User) Save(db *database.Database) error {
	user, err := db.UserCollection.InsertOne(db.Ctx, u)
	if err != nil {
		return errors.New("failed to save user: " + err.Error())

	}
	u.ID = user.InsertedID.(primitive.ObjectID)
	return nil
}

func (u *User) ComparePassword(expectedPassword string) bool {
	match, err := utils.ComparePasswordAndHash(expectedPassword, *u.Password)
	if err != nil {
		return false
	}
	return match
}

func GetUserByContext(c *gin.Context) *User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}

	return user.(*User)
}
