package User

import (
	"mime/multipart"

	usermodel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
)

type UserServiceInterface interface {
	CheckUserEmailExist(email string) bool
	CheckUserNameExist(username string) bool
	UserRegister(*usermodel.SignUpInput) *usermodel.DBResponse
	DeleteUserAccount(userName, Password string) bool
	GetUserByEmail(email string) *usermodel.User
	GetUserByUserName(username string) *usermodel.User
	UpdateUserProfile(userName, userProfile, bio string) (*usermodel.User, error)
	ChangePassword(userName, newPassword string) (*usermodel.User, error)
	UploadProfile(multipart.File,string) bool
}
