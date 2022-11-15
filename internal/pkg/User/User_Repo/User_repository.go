package userrepo

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"strings"

	userModel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
	"github.com/DavG20/Negarit_API/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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
		fmt.Println("error while decoding , in user_repo , means no user found ,")
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
		log.Println("error user not found")
		return nil
	}
	return user
}
func (userRepo *UserRepo) GetUserByUserName(userName string) (user *userModel.User) {
	user, err := userRepo.CheckUserNameExist(userName)
	if err != nil {
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

func (userRepo *UserRepo) DeleteUserAccount(userName, passwords string) (err error) {

	user, err := userRepo.CheckUserNameExist(userName)
	if err != nil {
		fmt.Println("eror in repo can't find user")
		return err
	}

	fmt.Println(user.Password, "password")
	passCompare := entity.ComparePasswordHash(user.Password, passwords)
	fmt.Println("pass here")
	if !passCompare {
		fmt.Println("eror pass comapare")
		return errors.New("eror comparing password")
	}
	fmt.Println("here tooo")
	filter := bson.D{{Key: "username", Value: user.Username}}
	_, err = userRepo.DB.Collection(entity.User).DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println("erorr in line 108")
		return err
	}
	return nil

}

func (userRepo *UserRepo) UpdateUserProfile(userName, userProfile, bio string) (user *userModel.DBResponse, err error) {
	filter := bson.D{{Key: "username", Value: userName}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "userprofile", Value: userProfile}, {Key: "bio", Value: bio}}}}
	_, err = userRepo.DB.Collection(entity.User).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	users := userModel.User{}
	err = userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(&users)
	user = userRepo.GetSecuredUser(&users)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (userRepo *UserRepo) ChangePassword(userName, newPassword string) (user *userModel.DBResponse, err error) {
	filter := bson.D{{Key: "username", Value: userName}}
	newhashedPass, err := entity.PasswordHash(newPassword)
	if err != nil {
		fmt.Println("eror password hashing")
		return nil, err
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: newhashedPass}}}}
	_, err = userRepo.DB.Collection(entity.User).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	users := userModel.User{}
	err = userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(&users)
	if err != nil {
		return nil, err
	}
	user = userRepo.GetSecuredUser(&users)
	return user, nil

}

func (userRepo UserRepo) UploadProfile(file multipart.File, fileName string) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("error while reading")
		return err

	}

	bucket, err := gridfs.NewBucket(userRepo.DB)
	if err != nil {
		fmt.Println("error creating bucket")
		return err
	}
	UploadStream, err := bucket.OpenUploadStream(fileName)
	if err != nil {
		fmt.Println("error creating writable stream")
		return err
	}
	fileSize, err := UploadStream.Write(data)
	if err != nil {
		fmt.Println("error writing data to stream")
		return err
	}
	fmt.Printf("the uploaded file name is %s , and the size is %d \n", fileName, fileSize)

	imageName := (strings.Split(fileName, "/"))
	userName := (strings.Split(imageName[len(imageName)-1], "."))[0]
    fmt.Printf("userName is %s,\n" ,userName)
	fileID := UploadStream.FileID
	filter := bson.D{{Key: "_id", Value: fileID}}
	contenttype := bson.D{{Key: "$set", Value: bson.D{{Key: "_id", Value: userName}, {Key: "contentType", Value: "image/jpeg"}}}}
	files := userRepo.DB.Collection("fs.files")
	updateResult, err := files.UpdateOne(context.TODO(), filter, contenttype)
	if err != nil {
		fmt.Println("error while updating file as Username")
		return err
	}
	fmt.Println(updateResult)
	return nil

}
