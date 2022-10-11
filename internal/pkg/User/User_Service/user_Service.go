package userservice

import (
	user "github.com/DavG20/Negarit_API/internal/pkg/User"
	userModel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
)


type UserService struct {
	UserRepoI user.UserRepoInterface
}

func NewUserService(repo user.UserRepoInterface) *UserService {
	return &UserService{UserRepoI: repo}
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

func (userService UserService) UserLogin(userInput *userModel.SignInInput) *userModel.DBResponse {
	return userService.UserRepoI.UserLogin(userInput)
}
