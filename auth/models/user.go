package models

import (
	"auth/db"
	"auth/services"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"-" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type UserDocument struct {
	ID   string `json:"id" bson:"_id"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"-" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func CreateUser(u User) (UserDocument, error) {
	collection := db.GetDefaultCollection()

	userDoc, err := findUserByEmail(u.Email)

	if err == nil {
		return UserDocument{}, errors.New("email is already taken")
	}

	u.CreatedAt = time.Now()
	u.Password, _ = services.ToHash(u.Password)

	id, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		return UserDocument{}, errors.New("error creating user")
	}

	userDoc, _ = findUserById(id.InsertedID.(primitive.ObjectID))

	return userDoc, nil
}

func AuthUser(u User) (UserDocument, error) {
	userDoc, err := findUserByEmail(u.Email)

	if err != nil {
		return UserDocument{}, errors.New("invalid credentials")
	}

	res := services.Compare(userDoc.Password, u.Password)

	if !res {
		return UserDocument{}, errors.New("invalid credentials")
	}

	return userDoc, nil
}

func findUserById(id primitive.ObjectID) (UserDocument, error) {
	return findUser(bson.D{{"_id", id}})
}

func findUserByEmail(email string) (UserDocument, error) {
	return findUser(bson.D{{"email", email}})
}

func findUser(filter bson.D) (UserDocument, error) {
	collection := db.GetDefaultCollection()

	user := UserDocument{}

	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}
