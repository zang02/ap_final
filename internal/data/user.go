package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserModel struct {
	DB *mongo.Database
}

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Login      string             `json:"login"`
	Email      string             `json:"email"`
	Name       string             `json:"fullname"`
	Password   string             `json:"passwordHash"`
	CreateDate string             `json:"create_date"`
}

func (u *UserModel) Insert(user User) error {
	user.CreateDate = humanDate(time.Now().Add(time.Hour * 6))

	_, err := u.DB.Collection("users").InsertOne(context.TODO(), user)

	return err
}

func (u *UserModel) GetByLogin(login string) (User, error) {
	var user User
	err := u.DB.Collection("users").FindOne(context.TODO(), bson.M{"login": login}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
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
