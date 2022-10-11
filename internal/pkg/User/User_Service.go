package User

import (
	usermodel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
)

type UserServiceInterface interface {
	CheckUserEmailExist(email string) bool
	CheckUserNameExist(username string) bool
	UserRegister(*usermodel.SignUpInput) *usermodel.DBResponse
}
