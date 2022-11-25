package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	friends "github.com/DavG20/Negarit_API/internal/pkg/Friends"
	session "github.com/DavG20/Negarit_API/internal/pkg/Session"
	userrepo "github.com/DavG20/Negarit_API/internal/pkg/User/User_Repo"
	user_service "github.com/DavG20/Negarit_API/internal/pkg/User/User_Service"
	DB "github.com/DavG20/Negarit_API/internal/pkg/db"
	"github.com/DavG20/Negarit_API/pkg/entity"

	apihandler "github.com/DavG20/Negarit_API/api/API_Handler"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

var once sync.Once
var friend *friends.Friends

func StartUp() {
	once.Do(
		func() {
			db = DB.ConnectMongoDB()
			if db == nil {
				log.Fatal("exiting...")
				os.Exit(1)
			}
			return
		},
	)
}

func init() {
	StartUp()
}

func main() {
	if db == nil {
		fmt.Println("envalid db")
	}

	userRepo := userrepo.NewUserRepo(db)
	userservice := user_service.NewUserService(userRepo)
	cookieHandler := session.NewCookieHanedler()
	userHandler := apihandler.NewUserHandler(*cookieHandler, userservice)
	fmt.Println(userHandler)

	friendsRepo := friends.NewFriendsRepo(db)
	friendsService := friends.NewFriendsService(friendsRepo)
	friendsHandler := apihandler.NewFriendsHandler(*cookieHandler, friendsService, userservice)

	fmt.Println("surver running ...")
	fr, _ := friendsService.GetFriendByUserName("DavG220", "DavG2000")
	deleted := friendsService.BlockFriend("DavG220", fr)
	frr, _ := friendsService.GetFriendByUserName("DavG220", "DavG2000")
	fmt.Println(deleted, "after block", frr)

	http.HandleFunc("/user/", userHandler.RegisterUser)
	http.HandleFunc("/user/login", (userHandler.UserLogin))
	http.HandleFunc("/user/logout", userHandler.UserLogout)
	http.HandleFunc("/user/deleteaccount", userHandler.DeleteUserAccount)
	http.HandleFunc("/user/updateprofile", userHandler.UpdateUserProfile)
	http.HandleFunc("/user/changepassword", userHandler.ChangePassword)
	http.HandleFunc("/user/userprofile", userHandler.MyProfile)
	http.HandleFunc("/user/searchuser", userHandler.SearchUser)
	http.HandleFunc("/user/uploadprofile", userHandler.UploadProfilePic)
	http.HandleFunc("/user/createfriends", friendsHandler.CreateFriendshipHandler)
	http.HandleFunc("/user/friend/delete", friendsHandler.DeleteFriendsHandler)
	http.HandleFunc("/", dispaly)
	http.HandleFunc("/del", delete)

	http.ListenAndServe(":8080", nil)

}

func test(res http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		res.Write([]byte("eror"))
		return
	}
	user := r.FormValue("username")
	fmt.Println(user)
	res.Write([]byte(user))
}

func dispaly(w http.ResponseWriter, r *http.Request) {

	filter := bson.D{{}}
	user := friends.Friends{}
	users := []friends.Friends{}
	cursor, err := db.Collection(entity.Friends).Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("no user found")
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&user)
		if err != nil {
			fmt.Println("error")
		}

		users = append(users, user)

	}
	w.Write(entity.MarshalIndentHelper(users))
}

func delete(res http.ResponseWriter, r *http.Request) {
	// filter := bson.D{{}}
	err := db.Collection(entity.Friends).Drop(context.TODO())
	if err != nil {
		res.Write([]byte("can't delete "))
		return
	}
	fmt.Println("the number of rows deleted is")
	res.Write([]byte("success"))
}
