package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

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

func main() {
	if db == nil {
		fmt.Println("envalid db")
	}

	defer db.Client().Disconnect(context.TODO())

}
