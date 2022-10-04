package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	user "github.com/DavG20/Negarit_API/internal/pkg/User"
	DB "github.com/DavG20/Negarit_API/internal/pkg/db"
	"github.com/DavG20/Negarit_API/internal/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	session "github.com/DavG20/Negarit_API/internal/pkg/Session"
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

func creatUser(w http.ResponseWriter, r *http.Request) {
	var userRepo user.UserRepo = user.UserRepo{DB: db}
	var user *user.SignUpInput

	err := json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user, "user")
	if err != nil {
		log.Fatal("error creating user")
	}
	users, err := userRepo.UserRegister(user)
	if err != nil {
		log.Fatal("error returning", err)
	}
	usr, err := json.MarshalIndent(users, "", "\t")
	w.Write(usr)
	fmt.Println(users)
}

func getUser(response http.ResponseWriter, r *http.Request) {
	var userRepo user.UserRepo = user.UserRepo{DB: db}
	var user user.User
	filter := bson.D{{"email", "dawit@gmail.com"}}

	err := userRepo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		response.Write([]byte("error while decoding or no user found"))
	}
	http.SetCookie(response, &http.Cookie{Name: "dav", Value: "davfuck u"})
	usr, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		response.Write([]byte("error while marshal indent"))
	}
	response.Write(usr)

}

func testCookie(w http.ResponseWriter, r *http.Request) {
	var cookieHandler session.CookieHandler = session.CookieHandler{}
	// var userInput user.SignInInput
	var session *session.Session
	// err := json.NewDecoder(r.Body).Decode(&userInput)
	// if err != nil {
	// 	w.Write([]byte("error while decoding"))
	// }

	// session := &session.Session{
	// 	Email: userInput.Email,
	// }

	cookie, err := cookieHandler.GetCookie(session)
	if err != nil {
		w.Write([]byte("error geting session"))
	}
	// fmt.Println(session.Email)
	http.SetCookie(w, &cookie)
	// usr, err := json.MarshalIndent(userInput, "", "\t")
	w.Header().Set("token", cookie.Value)
	w.Write([]byte("sample output"))

}

func testCookieValidate(w http.ResponseWriter, r *http.Request) {
	cookiehandler := session.CookieHandler{}
	fmt.Println(cookiehandler.ValidateCookie(r))

}
func logout(w http.ResponseWriter, r *http.Request) {
	cookieHandler := session.CookieHandler{}
	cookie, err := cookieHandler.RemoverCookie()
	if err != nil {
		http.SetCookie(w, nil)
		fmt.Println("error")
		w.Write([]byte("unauthorized user"))
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("successfuly logout"))
}
func main() {
	if db == nil {
		fmt.Println("envalid db")
	}

	defer db.Client().Disconnect(context.TODO())

	http.HandleFunc("/", testCookie)
	http.HandleFunc("/test", testCookieValidate)
	http.HandleFunc("/logout", logout)
	fmt.Println("surving...")
	http.ListenAndServe(":8080", nil)

}
