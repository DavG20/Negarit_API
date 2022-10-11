package userrepo

import (
	"context"
	"log"

	userModel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
	"github.com/DavG20/Negarit_API/internal/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	DB *mongo.Database
}

func newUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		DB: db,
	}

}
func (userRepo *UserRepo) CheckUserEmailExist(email string) (user *userModel.DBResponse, err error) {
	filter := bson.D{{Key: "email", Value: email}}
	err = userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserRepo) CheckUserNameExist(username string) (user *userModel.DBResponse, err error) {
	filter := bson.D{{Key: "username", Value: username}}

	err = userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (userRepo *UserRepo) UserRegister(userInput *userModel.SignUpInput) (user *userModel.DBResponse, err error) {
	//    i have to fix this
	_, err = userRepo.DB.Collection(entity.User).InsertOne(context.TODO(), userInput)
	if err != nil {
		return user, err
	}
	return user, nil

}

func (userRepo *UserRepo) GetUserByEmail(email string) (user *userModel.DBResponse) {
	user, err := userRepo.CheckUserEmailExist(email)
	if err != nil {
		log.Panicln("error user not found")
		return nil
	}
	return user
}

func (userRepo *UserRepo) UserLogin(userInput *userModel.SignInInput) (user *userModel.DBResponse) {
	pass, err := entity.PasswordHash(userInput.Password)
	if err != nil {
		log.Panicln("error while password hash")
	}
	filter := bson.D{{Key: "email", Value: userInput.Email}, {Key: "password", Value: pass}}
	err = userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		log.Panicln("error finding loging")
		return nil
	}
	return user
}
