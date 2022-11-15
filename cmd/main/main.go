package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	session "github.com/DavG20/Negarit_API/internal/pkg/Session"
	userModel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
	userrepo "github.com/DavG20/Negarit_API/internal/pkg/User/User_Repo"
	userservice "github.com/DavG20/Negarit_API/internal/pkg/User/User_Service"
	DB "github.com/DavG20/Negarit_API/internal/pkg/db"
	"github.com/DavG20/Negarit_API/pkg/entity"

	apihandler "github.com/DavG20/Negarit_API/api/API_Handler"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

var once sync.Once

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
	userservice := userservice.NewUserService(userRepo)
	cookieHandler := session.NewCookieHanedler()
	userHandler := apihandler.NewUserHandler(*cookieHandler, userservice)
	fmt.Println(userHandler)

	fmt.Println("surver running ...")
	print(entity.GenerateRandomString())

	http.HandleFunc("/user/", userHandler.RegisterUser)
	http.HandleFunc("/user/login", (userHandler.UserLogin))
	http.HandleFunc("/user/logout", userHandler.UserLogout)
	http.HandleFunc("/user/deleteaccount", userHandler.DeleteUserAccount)
	http.HandleFunc("/user/updateprofile", userHandler.UpdateUserProfile)
	http.HandleFunc("/user/changepassword", userHandler.ChangePassword)
	http.HandleFunc("/user/userprofile", userHandler.MyProfile)
	http.HandleFunc("/user/searchuser", userHandler.SearchUser)
	http.HandleFunc("/user/uploadprofile", userHandler.UploadProfilePic)
	

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
	// session, res := session.NewCookieHanedler().ValidateCookie(r)
	// if !res {
	// 	fmt.Println("can't get cookies")
	// }
	// fmt.Println(session.UserName, "h")

	filter := bson.D{{}}
	user := userModel.User{}
	users := []userModel.User{}
	cursor, err := db.Collection(entity.User).Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("no user found")
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&user)
		if err != nil {
			fmt.Println("error")
		}
		users = append(users, user)
		w.Write(entity.MarshalIndentHelper(users))

	}
}
