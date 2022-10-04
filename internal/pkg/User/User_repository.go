package user

import (
	"context"
	"log"

	"github.com/DavG20/Negarit_API/internal/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepoInterface interface {
	CheckUserEmail(email string) (*User, error)
	UserRegister(userInput SignUpInput) (*User, error)
	GetUserByEmail(email string) *User
	UserLogin(email, password string) *User
}

type UserRepo struct {
	DB *mongo.Database
}

func newUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		DB: db,
	}

}
func (userRepo *UserRepo) CheckUserEmail(email string) (user *User, err error) {
	filter := bson.D{{Key: "email", Value: email}}
	err = userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserRepo) UserRegister(userInput *SignUpInput) (user *User, err error) {
	//    i have to fix this
	_, err = userRepo.DB.Collection(entity.User).InsertOne(context.TODO(), userInput)
	if err != nil {
		return user, err
	}
	return user, nil

}

func (userRepo *UserRepo) GetUserByEmail(email string) (user *User) {
	user, err := userRepo.CheckUserEmail(email)
	if err != nil {
		log.Panicln("error user not found")
		return nil
	}
	return user
}

func (userRepo *UserRepo) UserLogin(email, password string) (user *User) {
	pass, err := entity.PasswordHash(password)
	if err != nil {
		log.Panicln("error while password hash")
	}
	filter := bson.D{{Key: "email", Value: email}, {Key: "password", Value: pass}}
	err = userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		log.Panicln("error finding loging")
		return nil
	}
	return user
}
