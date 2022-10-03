package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DavG20/Negarit_API/internal/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (userRepo *UserRepo) RegisterUser(inputuser *SignUpInput) (user *User, errs error) {
	user = &User{
		Username:    inputuser.Username,
		Email:       inputuser.Email,
		Password:    inputuser.Password,
		Userprofile: inputuser.Userprofile,
		Bio:         inputuser.Bio,
		CreatedAt:   time.Now(),
	}
	if userRepo.CheckUserEmailExist(inputuser.Email) {
		return nil, errors.New("email already registerd")
	}
	res, err := userRepo.DB.Collection(entity.User).InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Println("error inserting ")
		return nil, err
	}

	userId := entity.GetIdFromInsertedObjectId(res.InsertedID.(primitive.ObjectID))
	fmt.Println(userId, "this is userid")
	filter := bson.D{{"_id", user.Email}}
	update := bson.D{{"$set", bson.D{{"_id", userId}}}}
	// var u User
	resupd, err := userRepo.DB.Collection(entity.User).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("error updating userid")
	}
	fmt.Println("get user", user, "and ", resupd)
	return user, nil
}

func (userRepo *UserRepo) CheckUserEmailExist(email string) bool {
	filter := bson.D{{"email", email}}
	var user User
	err := userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return false
	}
	return true
}
