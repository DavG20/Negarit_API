package User

import (
	"mime/multipart"

	usermodel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
)

type UserRepoInterface interface {
	CheckUserEmailExist(email string) (*usermodel.User, error)
	CheckUserNameExist(username string) (*usermodel.User, error)
	UserRegister(userInput *usermodel.SignUpInput) (*usermodel.DBResponse, error)
	GetUserByEmail(email string) *usermodel.User
	GetUserByUserName(userName string) *usermodel.User
	GetSecuredUser(user *usermodel.User) *usermodel.DBResponse
	DeleteUserAccount(userName, Password string) error
	UpdateUserProfile(userName, userProfile, bio string) (*usermodel.DBResponse, error)
	ChangePassword(userName, newPassword string) (*usermodel.DBResponse, error)
	UploadProfile(multipart.File,string)error
}
