package db

import (
	"context"
	"fmt"
	"log"

	"github.com/DavG20/Negarit_API/pkg/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() *mongo.Database {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error to connect mongo db")
		return nil

	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("error pinging")
		return nil
	}
	fmt.Println("connecting mongo db ...")
	return client.Database(entity.DBName)

}
