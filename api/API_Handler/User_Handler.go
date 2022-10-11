package apihandler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	session "github.com/DavG20/Negarit_API/internal/pkg/Session"
	userModel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
	userservice "github.com/DavG20/Negarit_API/internal/pkg/User/User_Service"
	"github.com/DavG20/Negarit_API/internal/pkg/entity"
)

type SessionHandlerInterface interface {
	GetSessionHandler() *session.CookieHandler
}
type UserHandler struct {
	CookieHandler *session.CookieHandler
	UserServ      userservice.UserService
}

func NewUserHandler(cookieHandler session.CookieHandler, userServ userservice.UserService) *UserHandler {
	return &UserHandler{
		CookieHandler: &cookieHandler,
		UserServ:      userServ,
	}
}

func (userHandler UserHandler) getSessionHandler() session.CookieHandler {
	return *userHandler.CookieHandler
}

func (userHandler *UserHandler) UserLogin(response http.ResponseWriter, request *http.Request) {
	// User inpu in userModel
	var userInput *userModel.SignInInput
	// Resoponse incase the request will not success
	dbResponse := userModel.DBResponseFailed{}
	response.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(request.Body).Decode(&userInput) //decod the incoming request
	if err != nil {
		log.Println("invalid input")
		dbResponse.Message = "Invalid input, Please Try again!"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return

	}
	// check if user is registered or not

	userEmailExist := userHandler.UserServ.CheckUserEmailExist(userInput.Email)
	if userEmailExist {
		loginUser := userHandler.UserServ.GetUserByEmail(userInput.Email)
		if loginUser == nil {
			dbResponse.Message = "user doesn't exist"
			response.Write(entity.MarshalIndentHelper(dbResponse))
			return
		}
		// check is the password is correct
		ispasswordRight := entity.ComparePasswordHash(loginUser.Password, userInput.Password)
		if !ispasswordRight {
			dbResponse.Message = "password is not correct!"
			response.Write(entity.MarshalIndentHelper(dbResponse))
			return
		}
		// if the password is correct that means user is authenticated
		// so lets get dbresponse which is filtered for user's privacy

		userResponse := userHandler.UserServ.GetSecuredUser(loginUser)
		if userResponse == nil {
			dbResponse.Message = "Internal server error ,, plz try again"
			response.Write(entity.MarshalIndentHelper(dbResponse))
			return
		}
		session := session.Session{
			UserName: userResponse.Username,
		}
		cookies, err := userHandler.CookieHandler.GetCookie(&session)
		if err != nil {
			log.Println("errror in login 78")
			dbResponse.Message = "internal system error"
			response.Write(entity.MarshalIndentHelper(dbResponse))
			return
		}
		http.SetCookie(response, &cookies)

		response.Write(entity.MarshalIndentHelper(userResponse))
	}

	dbResponse.Message = "user not found with this email"
	response.Write(entity.MarshalIndentHelper(dbResponse))

}

func (userHandler UserHandler) RegisterUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	// the response struct from db or backend to frontend if the request is success
	user := userModel.SignUpInput{}
	// the response from backend if request is not successful
	userSignUpFailerResponse := userModel.DBResponseFailed{}
	// dbResponse:=userModel.DBResponse{}

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		log.Println("the input is not valid")
	}

	if !entity.ValidateUserName(user.Username) {
		userSignUpFailerResponse.Message = "invalid userName input"
		response.Write(entity.MarshalIndentHelper(userSignUpFailerResponse))
		return

	}

	if !entity.ValidatePassword(user.Password) {
		userSignUpFailerResponse.Message = "your password is not valid please try again"
		response.Write(entity.MarshalIndentHelper(userSignUpFailerResponse))
		return

	}

	if !entity.ValidateEmail(user.Email) {
		userSignUpFailerResponse.Message = "ur email is not valid"
		response.Write(entity.MarshalIndentHelper(userSignUpFailerResponse))
		return
	}

	// check if the email is registered
	emailExist := userHandler.UserServ.CheckUserEmailExist(user.Email)
	if emailExist {
		userSignUpFailerResponse.Message = "Email already registered! plz provide another email"
		response.Write(entity.MarshalIndentHelper(userSignUpFailerResponse))
		return
	}

	// check  username  if it already registered

	UserNameExist := userHandler.UserServ.CheckUserNameExist(user.Username)
	if UserNameExist {
		userSignUpFailerResponse.Message = "Username already Taken , try another one"
		response.Write(entity.MarshalIndentHelper(userSignUpFailerResponse))
		return
	}

	pass, err := entity.PasswordHash(user.Password)
	if err != nil {
		log.Println("error hashing  password,  User_Handler ")
		userSignUpFailerResponse.Message = "password hashing error"
		response.Write(entity.MarshalIndentHelper(userSignUpFailerResponse))
		return

	}

	user.CreatedOn = time.Now()
	user.Password = pass

	// call the service method to create user document , and will return DBUserResponse
	res := userHandler.UserServ.UserRegister(&user)
	if res == nil {
		response.WriteHeader(http.StatusInternalServerError)
		userSignUpFailerResponse.Message = "Problem in the server ,,,,, plz try again"
		response.Write(entity.MarshalIndentHelper(userSignUpFailerResponse))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(entity.MarshalIndentHelper(res))

}
