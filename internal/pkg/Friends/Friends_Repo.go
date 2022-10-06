package friends

import (
	"context"
	"log"

	message "github.com/DavG20/Negarit_API/internal/pkg/Message"
	"github.com/DavG20/Negarit_API/internal/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Friends_Repo_Interface interface {
	CreateFriendship(*message.Message) (*Friends, error)
}

type Friends_Repo struct {
	DB mongo.Database
}

func NewFriendsRepo(db *mongo.Database) *Friends_Repo {
	return &Friends_Repo{
		DB: *db,
	}
}

func (friends_Repo *Friends_Repo) CreateFriendship(message *message.Message) (friends *Friends, err error) {
	// res:=friends_Repo.AreTheyFriend(message)
	
	
	friends = &Friends{
		Friend_A_UserName: message.Sender_UserName,
		Friend_B_UserName: message.Receiver_UserName,
		Message:           append(friends.Message, message),
	}
	return friends, err

}

func (friends_Repo *Friends_Repo) AreTheyFriend(message *message.Message) bool {
	friends := &Friends{}
	filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "friend_a_username", Value: message.Sender_UserName}, {Key: "friend_b_username", Value: message.Receiver_UserName}}, bson.D{{Key: "friend_b_username", Value: message.Sender_UserName}, {Key: "friend_a_username", Value: message.Receiver_UserName}}}}}
	err := friends_Repo.DB.Collection(entity.Friends).FindOne(context.TODO(), filter).Decode(friends)
	if err != nil {
		log.Println("error finding friends")
		return false
	}
	return true
}
