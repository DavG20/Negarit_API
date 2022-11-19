package userservice

import (
	"fmt"
	"mime/multipart"

	user "github.com/DavG20/Negarit_API/internal/pkg/User"
	userModel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
)

type UserService struct {
	UserRepoI user.UserRepoInterface
}

func NewUserService(repo user.UserRepoInterface) UserService {
	return UserService{UserRepoI: repo}
}

func (userService *UserService) CheckUserEmailExist(email string) bool {
	_, err := userService.UserRepoI.CheckUserEmailExist(email)
	if err != nil {
		return false
	}
	return true
}

func (userService *UserService) CheckUserNameExist(username string) bool {
	_, err := userService.UserRepoI.CheckUserNameExist(username)
	if err != nil {
		return false
	}
	return true

}

func (userService *UserService) UserRegister(userInput *userModel.SignUpInput) *userModel.DBResponse {
	user, err := userService.UserRepoI.UserRegister(userInput)
	if err != nil {
		return nil
	}
	return user
}

func (userService *UserService) GetUserByEmail(email string) *userModel.User {
	return userService.UserRepoI.GetUserByEmail(email)
}
func (userService *UserService) GetUserByUserName(userName string) *userModel.User {
	return userService.UserRepoI.GetUserByUserName(userName)
}

func (userService *UserService) GetSecuredUser(user *userModel.User) *userModel.DBResponse {
	return userService.UserRepoI.GetSecuredUser(user)
}

func (userService *UserService) DeleteUserAccount(userName, Password string) bool {
	err := userService.UserRepoI.DeleteUserAccount(userName, Password)
	if err != nil {
		fmt.Println("error in service line 54")
		return false
	}
	return true
}

func (userService *UserService) UpdateUserProfile(userName, userProfile, bio string) (*userModel.DBResponse, bool) {
	user, err := userService.UserRepoI.UpdateUserProfile(userName, userProfile, bio)
	if err != nil {
		return nil, false
	}
	return user, true
}

func (userService *UserService) ChangePassword(userName, newPassword string) (*userModel.DBResponse, bool) {
	user, err := userService.UserRepoI.ChangePassword(userName, newPassword)
	if err != nil {
		return nil, false
	}
	return user, true
}

func (userService *UserService) UploadProfile(file multipart.File, fileName string) bool {
	err := userService.UserRepoI.UploadProfile(file, fileName)
	if err != nil {
		return false
	}
	return true
}
