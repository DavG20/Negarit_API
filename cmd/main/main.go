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
	"github.com/DavG20/Negarit_API/internal/pkg/entity"

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
	http.HandleFunc("/", userHandler.UserLogin)
	http.ListenAndServe(":8080", nil)

}
func dispaly(w http.ResponseWriter, r *http.Request) {
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
