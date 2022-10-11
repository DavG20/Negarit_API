package User

import (
	usermodel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
)

type UserRepoInterface interface {
	CheckUserEmailExist(email string) (*usermodel.User, error)
	CheckUserNameExist(username string) (*usermodel.User, error)
	UserRegister(userInput *usermodel.SignUpInput) (*usermodel.DBResponse, error)
	GetUserByEmail(email string) *usermodel.User
	CheckUserLogin(userInput *usermodel.SignInInput) *usermodel.DBResponse
	GetSecuredUser(user *usermodel.User) *usermodel.DBResponse
}
