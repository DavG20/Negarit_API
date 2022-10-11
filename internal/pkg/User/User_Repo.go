package User

import (
	usermodel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
)

type UserRepoInterface interface {
	CheckUserEmailExist(email string) (*usermodel.DBResponse, error)
	CheckUserNameExist(username string) (*usermodel.DBResponse, error)
	UserRegister(userInput *usermodel.SignUpInput) (*usermodel.DBResponse, error)
	GetUserByEmail(email string) *usermodel.DBResponse
	UserLogin(userInput *usermodel.SignInInput) *usermodel.DBResponse
}
