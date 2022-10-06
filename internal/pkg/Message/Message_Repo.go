package message

import "go.mongodb.org/mongo-driver/mongo"

// mesage repo struct
type MessageRepo struct {
	DB mongo.Database
}

// to create new message repo
func NewMessageRepo(db mongo.Database) *MessageRepo {
	return &MessageRepo{
		DB: db,
	}
}

