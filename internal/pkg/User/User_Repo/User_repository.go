package userrepo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	userModel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
	"github.com/DavG20/Negarit_API/internal/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	DB *mongo.Database
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		DB: db,
	}

}
func (userRepo *UserRepo) CheckUserEmailExist(email string) (user *userModel.User, err error) {
	filter := bson.D{{Key: "email", Value: email}}
	fmt.Println("Ok")
	errs := userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(&user)

	if errs != nil {
		fmt.Println("error while decoding")
		return nil, errs
	}

	return user, nil
}

func (userRepo *UserRepo) CheckUserNameExist(username string) (user *userModel.User, err error) {
	filter := bson.D{{Key: "username", Value: username}}

	err = userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (userRepo *UserRepo) UserRegister(userInput *userModel.SignUpInput) (user *userModel.DBResponse, err error) {
	//    i have to fix this
	res, err := userRepo.DB.Collection(entity.User).InsertOne(context.TODO(), userInput)
	if err != nil {
		return nil, err
	}
	userId := entity.GetIdFromInsertedObjectId(res.InsertedID.(primitive.ObjectID))
	user = &userModel.DBResponse{
		UserId:      userId,
		Email:       userInput.Email,
		Username:    userInput.Username,
		Bio:         userInput.Bio,
		Userprofile: userInput.Userprofile,
		CreatedOn:   userInput.CreatedOn,
	}
	return user, nil

}

func (userRepo *UserRepo) GetUserByEmail(email string) (user *userModel.User) {
	user, err := userRepo.CheckUserEmailExist(email)
	if err != nil {
		log.Panicln("error user not found")
		return nil
	}
	return user
}

func (userRepo *UserRepo) GetSecuredUser(user *userModel.User) *userModel.DBResponse {
	return &userModel.DBResponse{
		UserId:      user.UserId,
		Email:       user.Email,
		Username:    user.Username,
		Userprofile: user.Userprofile,
		Bio:         user.Bio,
		CreatedOn:   user.CreatedOn,
	}
}

func (userRepo *UserRepo) CheckUserLogin(userInput *userModel.SignInInput) (DBRuser *userModel.DBResponse) {
	pass, err := entity.PasswordHash(userInput.Password)
	if err != nil {
		log.Panicln("error while password hash")
	}
	filter := bson.D{{Key: "email", Value: userInput.Email}, {Key: "password", Value: pass}}
	err = userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(&DBRuser)
	if err != nil {
		log.Println("error  logging  ", pass)
		return nil
	}
	return DBRuser
}

func (userRepo *UserRepo) CheckLogout(request *http.Request) bool {
	return true

}
