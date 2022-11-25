package apihandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	session "github.com/DavG20/Negarit_API/internal/pkg/Session"
	userModel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
	usermodel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
	userservice "github.com/DavG20/Negarit_API/internal/pkg/User/User_Service"
	"github.com/DavG20/Negarit_API/pkg/entity"
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

// check if user is login first to get some service
func (userHadler *UserHandler) CheckUserLogin(request *http.Request) (*session.Session, bool) {
	session, isValid := userHadler.CookieHandler.ValidateCookie(request)
	if isValid {
		return session, isValid
	}
	return nil, isValid

}

// check if user if logged out or not

func (userHandler *UserHandler) Authenticated(endpoint http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, isLogin := userHandler.CheckUserLogin(r)
		if !isLogin {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Println("U r not logged In , plz login first")
			return
		}
		if session == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return

		}
		endpoint(w, r)

	})
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

	if userInput.Email == "" {
		dbResponse.Message = "empty email , please provide ur email"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	if userInput.Password == "" {
		dbResponse.Message = "empty password , please provide ur password"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}

	// check if user is registered or not

	userEmailExist := userHandler.UserServ.CheckUserEmailExist(userInput.Email)
	if userEmailExist {
		loginUser := userHandler.UserServ.GetUserByEmail(userInput.Email)
		fmt.Println("here  in login func line 82")
		if loginUser == nil {
			dbResponse.Message = "user doesn't exist"
			response.Write(entity.MarshalIndentHelper(dbResponse))
			return
		}
		fmt.Println(loginUser.Password, "password")
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
			log.Println("errror in login 78 , can;t get session")
			dbResponse.Message = "internal system error"
			response.Write(entity.MarshalIndentHelper(dbResponse))
			return
		}
		http.SetCookie(response, &cookies)
		response.Write(entity.MarshalIndentHelper(userResponse))
		return
	}

	dbResponse.Message = "user not found with this email"
	response.Write(entity.MarshalIndentHelper(dbResponse))

}

func (userHandler *UserHandler) UserLogout(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	session, isValid := userHandler.CookieHandler.ValidateCookie(request)
	dbResponseFailed := userModel.DBResponseFailed{}
	if !isValid {
		dbResponseFailed.Message = "U r not logged in, u r logged out already "
		response.Write(entity.MarshalIndentHelper(dbResponseFailed))
		return

	}
	if session == nil {
		dbResponseFailed.Message = "unKnown user"
		response.Write(entity.MarshalIndentHelper(dbResponseFailed))
		return
	}
	userExist := userHandler.UserServ.CheckUserNameExist(session.UserName)
	if !userExist {
		dbResponseFailed.Message = "User not found!, u r trying to logout with out registered"
		response.Write(entity.MarshalIndentHelper(dbResponseFailed))
		return
	}
	logoutCookie, err := userHandler.CookieHandler.RemoveCookie()
	if err != nil {
		dbResponseFailed.Message = "error logout plz try again!"
		response.Write(entity.MarshalIndentHelper(dbResponseFailed))
		return
	}
	//    i have to edit something in here
	http.SetCookie(response, &logoutCookie)
	dbResponseFailed.Message = "logged Out successfuly"
	response.WriteHeader(http.StatusOK)
	response.Write(entity.MarshalIndentHelper(dbResponseFailed))

}

func (userHandler UserHandler) RegisterUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	// the response struct from db or backend to frontend if the request is success
	user := userModel.SignUpInput{}
	// the response from backend if request is not successful
	dbResponse := userModel.DBResponseFailed{}
	// dbResponse:=userModel.DBResponse{}

	err := json.NewDecoder(request.Body).Decode(&user)
	fmt.Println(user, "user")
	if err != nil {
		log.Println("the input is not valid")
	}

	if !entity.ValidateUserName(user.Username) {
		dbResponse.Message = "invalid userName input"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return

	}

	if !entity.ValidatePassword(user.Password) {
		dbResponse.Message = "your password is not valid please try again"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return

	}

	if !entity.ValidateEmail(user.Email) {
		dbResponse.Message = "Ur email is not valid"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}

	// check if the email is registered
	emailExist := userHandler.UserServ.CheckUserEmailExist(user.Email)
	if emailExist {
		dbResponse.Message = "Email already registered! plz provide another email"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}

	// check  username  if it already registered

	UserNameExist := userHandler.UserServ.CheckUserNameExist(user.Username)
	if UserNameExist {
		dbResponse.Message = "Username already Taken , try another one"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}

	pass, err := entity.PasswordHash(user.Password)
	if err != nil {
		log.Println("error hashing  password,  User_Handler ")
		dbResponse.Message = "password hashing error"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return

	}

	user.CreatedOn = time.Now()
	user.Password = pass

	// call the service method to create user document , and will return DBUserResponse
	res := userHandler.UserServ.UserRegister(&user)
	if res == nil {
		response.WriteHeader(http.StatusInternalServerError)
		dbResponse.Message = "Problem in the server ,,,,, plz try again"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	//    if the user is registered successfully he/she doesn't has to login again
	// let the user login after registration , to do that i have to send a cookie/session for them for the future use

	session := session.Session{
		UserName: user.Username,
	}
	cookie, err := userHandler.CookieHandler.GetCookie(&session)
	if err != nil {
		log.Println("regsiration cookie can't set, line 171")
		dbResponse.Message = "You have been registered successfuly try to login"
		response.Write(entity.MarshalIndentHelper(dbResponse))
	}
	http.SetCookie(response, &cookie)
	response.WriteHeader(http.StatusOK)
	response.Write(entity.MarshalIndentHelper(res))

}

func (userHandler *UserHandler) DeleteUserAccount(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	type password struct {
		Pass string `json:"password"`
	}
	inputPass := password{}
	responseMessage := usermodel.DBResponseFailed{}
	err := json.NewDecoder(request.Body).Decode(&inputPass)
	if err != nil {
		responseMessage.Message = "invalid password , please try again"
		response.Write(entity.MarshalIndentHelper(responseMessage))
		return

	}
	userPassword := inputPass.Pass
	session, isValid := userHandler.CookieHandler.ValidateCookie(request)
	if !isValid {
		responseMessage.Message = "U have to login first to delete ur account"
		response.Write(entity.MarshalIndentHelper(responseMessage))
		return
	}
	fmt.Println(session.UserName, "userName")
	isUserExist := userHandler.UserServ.CheckUserNameExist(session.UserName)
	if !isUserExist {
		responseMessage.Message = "User Not found"
		response.Write(entity.MarshalIndentHelper(responseMessage))
		return
	}
	isDeleted := userHandler.UserServ.DeleteUserAccount(session.UserName, userPassword)
	fmt.Println("error not here")
	if !isDeleted {
		responseMessage.Message = "error while deleting "
		response.Write(entity.MarshalIndentHelper(responseMessage))
		return
	}
	clearSession, err := userHandler.CookieHandler.RemoveCookie()
	if err != nil {
		fmt.Println("error clearing session")
		return
	}
	http.SetCookie(response, &clearSession)
	responseMessage.Message = "successuly deleted"
	response.WriteHeader(http.StatusOK)
	response.Write(entity.MarshalIndentHelper(responseMessage))

}

// update user profile

func (userHandler *UserHandler) UpdateUserProfile(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userInput := userModel.UpdateUserProfile{}
	dbResonse := userModel.UserProfileUpdateResponse{}
	// check the user authentication
	session, isValid := userHandler.CookieHandler.ValidateCookie(request)

	if !isValid {
		dbResonse.Message = "Un authorized user"
		dbResonse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResonse))
		return
	}
	isuserExist := userHandler.UserServ.CheckUserNameExist(session.UserName)
	if !isuserExist {
		dbResonse.Message = "User Unknown , please try again"
		dbResonse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResonse))
		return
	}

	err := json.NewDecoder(request.Body).Decode(&userInput)
	if err != nil || userInput.Bio == "" || userInput.Userprofile == "" {
		dbResonse.Message = "Invalid Input"
		dbResonse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResonse))
		return
	}
	user := userHandler.UserServ.GetUserByUserName(session.UserName)
	fmt.Println(user)
	if user.Userprofile == userInput.Userprofile && user.Bio == userInput.Bio {
		dbResonse.Message = "ur input is the same as the old profile, please try to provide a new one"
		dbResonse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResonse))
		return

	}
	UpdatedUser, isUpdated := userHandler.UserServ.UpdateUserProfile(session.UserName, userInput.Userprofile, userInput.Bio)
	if !isUpdated {
		dbResonse.Message = "internal server error, can't update profile , please try again!"
		dbResonse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResonse))
		return
	}
	dbResonse.Message = "successfuly updated userProfile"
	dbResonse.Success = true
	dbResonse.User = append(dbResonse.User, *UpdatedUser)
	response.WriteHeader(http.StatusOK)
	response.Write(entity.MarshalIndentHelper(dbResonse))

}

// change password
func (userHandler *UserHandler) ChangePassword(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	dbresponse := userModel.UserProfileUpdateResponse{}
	userIpnut := userModel.UpdatePassword{}

	// check is the request is sent from authorized user
	session, isValid := userHandler.CookieHandler.ValidateCookie(request)
	if !isValid {
		dbresponse.Message = " not Authorized to change password"
		dbresponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbresponse))
		return
	}
	// check is the user is exist from the request session
	isUserexist := userHandler.UserServ.CheckUserNameExist(session.UserName)
	if !isUserexist {
		dbresponse.Success = false
		dbresponse.Message = "User not found "
		response.Write(entity.MarshalIndentHelper(dbresponse))
		return
	}
	// get user's info to compare with the incoming password
	user := userHandler.UserServ.GetUserByUserName(session.UserName)
	if user == nil {
		dbresponse.Message = "sorry can't find user"
		dbresponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbresponse))
		return
	}

	// encode the requset body with custom struct

	err := json.NewDecoder(request.Body).Decode(&userIpnut)
	if err != nil {
		dbresponse.Message = "Invalid input, please provide the correct information"
		dbresponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbresponse))
		return
	}

	// check is the input is not empty string
	if userIpnut.NewPassword == "" || userIpnut.OldPassword == "" {
		dbresponse.Message = "empty password plz fill according to the requirment"
		dbresponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbresponse))
		return
	}

	// check if the old password is correct
	isvalidpass := entity.ComparePasswordHash(user.Password, userIpnut.OldPassword)
	if !isvalidpass {
		dbresponse.Message = "Incorrect old password, "
		dbresponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbresponse))
		return
	}

	updatedUser, isSuccess := userHandler.UserServ.ChangePassword(session.UserName, userIpnut.NewPassword)
	if !isSuccess {
		dbresponse.Message = "can't change password, Internal server error plz try again"
		dbresponse.Success = isSuccess
		response.Write(entity.MarshalIndentHelper(dbresponse))
		return
	}

	dbresponse.Message = "password changed successfuly"
	dbresponse.Success = true
	dbresponse.User = append(dbresponse.User, *updatedUser)
	response.WriteHeader(http.StatusOK)
	response.Write(entity.MarshalIndentHelper(dbresponse))
}

func (userHandler *UserHandler) MyProfile(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	dbResponse := usermodel.UserProfile{}
	session, isValid := userHandler.CookieHandler.ValidateCookie(request)
	if !isValid {
		dbResponse.Message = "Unauthorized User"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return

	}
	if session == nil {
		dbResponse.Message = "Internal server error plz try again"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	isUserExist := userHandler.UserServ.CheckUserNameExist(session.UserName)
	if !isUserExist {
		dbResponse.Message = "User not found"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	user := userHandler.UserServ.GetUserByUserName(session.UserName)
	if user == nil {
		dbResponse.Message = "internal server error , plz try again"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}

	UserProfile := userHandler.UserServ.GetSecuredUser(user)
	if UserProfile == nil {
		dbResponse.Message = "profile load error please try again"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}

	response.WriteHeader(http.StatusOK)
	dbResponse.Message = "Profile load successfuly"
	dbResponse.Success = true
	dbResponse.UserPro = append(dbResponse.UserPro, *UserProfile)
	response.Write(entity.MarshalIndentHelper(dbResponse))

}

func (userHandler *UserHandler) SearchUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	dbResponse := usermodel.UserProfile{}
	userName := request.FormValue("username")
	session, isValid := userHandler.CookieHandler.ValidateCookie(request)
	if !isValid {
		dbResponse.Message = "Unauthorized user"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	if session == nil {
		dbResponse.Message = "Session expired plz login again and try"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	isUserExist := userHandler.UserServ.CheckUserNameExist(session.UserName)
	// incase the request is not lagal and a session username is modified
	if !isUserExist {
		dbResponse.Message = "you are not a legal client to get this service ,"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	if userName == "" {
		dbResponse.Message = "Empty field ,"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	isSearchedUserFound := userHandler.UserServ.CheckUserNameExist(userName)
	if !isSearchedUserFound {
		dbResponse.Message = "user not found by this username"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	user := userHandler.UserServ.GetUserByUserName(userName)
	if user == nil {
		dbResponse.Message = "internal server error, please try again"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	SearchedUser := userHandler.UserServ.GetSecuredUser(user)
	if SearchedUser == nil {
		dbResponse.Message = "internal problem , please try again"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	response.WriteHeader(http.StatusOK)
	dbResponse.Message = "User Loaded successfuly"
	dbResponse.Success = true
	dbResponse.UserPro = append(dbResponse.UserPro, *SearchedUser)
	response.Write(entity.MarshalIndentHelper(dbResponse))

}

func (userHandler *UserHandler) UploadProfilePic(response http.ResponseWriter, request *http.Request) {
	session, isvalid := userHandler.CookieHandler.ValidateCookie(request)
	dbResponse := userModel.UserProfile{}
	if !isvalid {
		dbResponse.Message = "un Uthorized user"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	if session == nil {
		dbResponse.Message = "internal server error"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	err := request.ParseForm()
	if err != nil {
		dbResponse.Message = "invalid input"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
	}
	file, header, err := request.FormFile("file")
	if err != nil {
		dbResponse.Message = "invalid file path"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}

	defer file.Close()

	isImageValid := entity.IsProfileValid(header.Filename)
	if !isImageValid {
		dbResponse.Message = "only image type is allowed!"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	ImageName := "assests/image/profile/" + session.UserName + "." + entity.Getextension(header.Filename)
	user := userHandler.UserServ.GetUserByUserName(session.UserName)
	if user == nil {
		response.WriteHeader(http.StatusInternalServerError)
		dbResponse.Message = "Internal Server error , please try again"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	fmt.Println(ImageName, "here")
	Success := userHandler.UserServ.UploadProfile(file, ImageName)
	if !Success {
		dbResponse.Message = "error stroring image . please try again"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	updateUser, err := userHandler.UserServ.UserRepoI.UpdateUserProfile(user.Username, ImageName, user.Bio)
	if err != nil {
		dbResponse.Message = "Internal error , please try again"
		dbResponse.Success = false
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}

	fmt.Println(updateUser)
}
