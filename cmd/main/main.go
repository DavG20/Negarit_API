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

func creatUser(w http.ResponseWriter, r *http.Request) {
	var userRepo user.UserRepo = user.UserRepo{DB: db}
	var user *user.SignUpInput

	err := json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user, "user")
	if err != nil {
		log.Fatal("error creating user")
	}
	users, err := userRepo.RegisterUser(user)
	if err != nil {
		log.Fatal("error returning", err)
	}
	usr, err := json.MarshalIndent(users, "", "\t")
	w.Write(usr)
	fmt.Println(users)
}
func main() {
	if db == nil {
		fmt.Println("envalid db")
	}

	// user,err:=userRepo.RegisterUser()

	defer db.Client().Disconnect(context.TODO())

	http.HandleFunc("/", creatUser)
	fmt.Println("surving...")
	http.ListenAndServe(":8080", nil)

}
