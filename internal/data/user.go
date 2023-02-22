package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserModel struct {
	DB *mongo.Database
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Login        string             `json:"login"`
	Email        string             `json:"email"`
	Fullname     string             `json:"fullname"`
	Passwordhash string             `json:"passwordhash"`
	CreateDate   string             `json:"create_date"`
}

func (u *UserModel) Insert(email string, login string, passwordhash string, fullname string) error {
	newUser := User{
		Email:        email,
		Passwordhash: passwordhash,
		Login:        login,
		Fullname:     fullname,
		CreateDate:   humanDate(time.Now().Add(time.Hour * 6)),
	}
	_, err := u.DB.Collection("users").InsertOne(context.TODO(), newUser)

	return err
}

// func (u *UserModel) Validate() error {
// 	u.DB.Collection("users").Aggregate()

// 	for _, user := range users {
// 		if user.Login == u.Login {
// 			return errors.New("User already exists")
// 		}
// 	}
// 	return nil
// }
