package apihandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	frinds "github.com/DavG20/Negarit_API/internal/pkg/Friends"
	message "github.com/DavG20/Negarit_API/internal/pkg/Message"
	session "github.com/DavG20/Negarit_API/internal/pkg/Session"
	usermodel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
	userservice "github.com/DavG20/Negarit_API/internal/pkg/User/User_Service"
	"github.com/DavG20/Negarit_API/pkg/entity"
)

type FriendsHandler struct {
	CookieHandler  *session.CookieHandler
	FriendsService frinds.FriendsService
	UserService    userservice.UserService
}

func NewFriendsHandler(cookieHandler session.CookieHandler, friendsService frinds.FriendsService, userService userservice.UserService) *FriendsHandler {
	return &FriendsHandler{
		CookieHandler:  &cookieHandler,
		FriendsService: friendsService,
		UserService:    userService,
	}
}

// Create Friendship handler for incoming request , the user input will be message's content and reciever's username
// we also use create friendship habdler when we want to append a new message to an existing friendship inbox
func (friendsHandler *FriendsHandler) CreateFriendshipHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "appliaction/json")
	dbResponse := usermodel.DBResponseFailed{}
	message := message.Message{}

	session, isLoggedIn := friendsHandler.CookieHandler.ValidateCookie(request)
	if !isLoggedIn {
		dbResponse.Message = "U are not logged in , login or register first befor u send friend request :)"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	fmt.Println(session)
	isUserExist := friendsHandler.UserService.CheckUserNameExist(session.UserName)
	if !isUserExist {
		dbResponse.Message = "U are not registered to send friend request"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return

	}

	// when user wants send message , but they are not friends yet, I have to create friendship for them in the first time
	err := json.NewDecoder(request.Body).Decode(&message)
	if err != nil {
		dbResponse.Message = "invalod user input please try again"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	isFriendExist := friendsHandler.UserService.CheckUserNameExist(message.Receiver_UserName)
	if !isFriendExist {
		dbResponse.Message = "friends username doesn't exist, "
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	if message.Receiver_UserName == session.UserName {
		dbResponse.Message = "Can't send self message , please try another reciever's userName"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}

	message.Sender_UserName = session.UserName
	areTheyFriends := friendsHandler.FriendsService.AreTheyFriend(message.Sender_UserName, message.Receiver_UserName)
	if areTheyFriends {
		fmt.Println("in here")
		friend, err := friendsHandler.FriendsService.GetFriendByUserName(message.Sender_UserName, message.Receiver_UserName)
		if err != nil {
			dbResponse.Message = "internal server error , please try again"
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(entity.MarshalIndentHelper(dbResponse))
			return
		}
		message.Message_Id = len(friend.Message)
		friend.Message = append(friend.Message, message)

		isAppendSuccess := friendsHandler.FriendsService.AppendMessage(friend)
		if !isAppendSuccess {
			dbResponse.Message = "can't send message to friend, internal problem"
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(entity.MarshalIndentHelper(dbResponse))
			return
		}
		response.WriteHeader(http.StatusOK)
		response.Write(entity.MarshalIndentHelper(friend))
		return
	}
	message.Message_Id = 0
	friend, err := friendsHandler.FriendsService.CreateFriendship(message)
	if err != nil {
		dbResponse.Message = "eror creating friendship or sending message"
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return

	}

	response.WriteHeader(http.StatusOK)
	response.Write(entity.MarshalIndentHelper(friend))

}

// delete friends by using his/her username
// the username will send through formvalue request
func (friendsHandler *FriendsHandler) DeleteFriendsHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	dbResponse := usermodel.DBResponseFailed{}
	// check is the request is authorized
	// and this also give me the username of the user who sent the request
	session, isValid := friendsHandler.CookieHandler.ValidateCookie(request)
	if !isValid {
		dbResponse.Message = "Unauthorized User to delet a friend, please login in first"
		response.WriteHeader(http.StatusUnauthorized)
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	// Check is user is in our system
	isUserExist := friendsHandler.UserService.CheckUserNameExist(session.UserName)
	if !isUserExist {
		dbResponse.Message = "U are not Registered , "
		response.WriteHeader(http.StatusUnauthorized)
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	// get the friends username and check is he/she is in our system
	friendsUserName := request.FormValue("friendsUserName")
	isFriendExist := friendsHandler.UserService.CheckUserNameExist(friendsUserName)
	if !isFriendExist {
		dbResponse.Message = "Friends Username doesn't exist, please provide another one"
		response.WriteHeader(http.StatusBadRequest)
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	// check if the user send his/her username incase :)
	if friendsUserName == session.UserName {
		dbResponse.Message = "you are not friend' with ur self:), got u"
		response.WriteHeader(http.StatusBadRequest)
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	// now both user's are in our system so lets check if they are friends with each other
	areTheyFriends := friendsHandler.FriendsService.AreTheyFriend(session.UserName, friendsUserName)
	if !areTheyFriends {
		dbResponse.Message = fmt.Sprintf("u r not friend with %s \n", friendsUserName)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}

	isDeleted := friendsHandler.FriendsService.DeleteFriends(session.UserName, friendsUserName)
	if !isDeleted {
		dbResponse.Message = "sorry some problem happens in our server , can't delete friend , plese try again"
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(entity.MarshalIndentHelper(dbResponse))
		return
	}
	dbResponse.Message = "friend deleted successfuly, "
	response.WriteHeader(http.StatusOK)
	response.Write(entity.MarshalIndentHelper(dbResponse))
}
