package friends

import (
	"context"
	"errors"
	"log"

	message "github.com/DavG20/Negarit_API/internal/pkg/Message"
	user_service "github.com/DavG20/Negarit_API/internal/pkg/User/User_Service"
	"github.com/DavG20/Negarit_API/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Friends_Repo_Interface interface {
	CreateFriendship(*message.Message) (*Friends, error)
	AreTheyFriend(*message.Message) bool
	GetFriendByUserName(ownerUserName, friendsUserName string) (*Friends, error)
}

type Friends_Repo struct {
	DB mongo.Database
}

func NewFriendsRepo(db *mongo.Database) *Friends_Repo {
	return &Friends_Repo{
		DB: *db,
	}
}

// create friendship function
// to create friendship from the message
func (friends_Repo *Friends_Repo) CreateFriendship(message *message.Message) (friends *Friends, err error) {
	userServ := user_service.UserService{}
	result1 := userServ.CheckUserNameExist(message.Sender_UserName)
	result2 := userServ.CheckUserNameExist(message.Receiver_UserName)
	if result1 && result2 {
		friends = &Friends{
			Friend_A_UserName: message.Sender_UserName,
			Friend_B_UserName: message.Receiver_UserName,
			Message:           append(friends.Message, message),
		}
		Friendship_id := ""
		res, err := friends_Repo.DB.Collection(entity.Friends).InsertOne(context.TODO(), friends)
		if err != nil {
			log.Println("error creating friendship")
			return nil, err

		}

		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			Friendship_id = entity.GetIdFromInsertedObjectId(oid)
		}
		friends.Friend_Id = Friendship_id
		return friends, nil
	}
	return nil, errors.New("The user is not found")

}

// check if they are friends when the  user sends message in the first time
// if they are friend return true if not return false

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

// to get friends by his/her username  and return friends
// they must be friends already

func (friends_Repo *Friends_Repo) GetFriendByUserName(ownerUserName, friendsUserName string) (friends *Friends, err error) {

	filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "friend_a_username", Value: ownerUserName}, {Key: "friends_b_username", Value: friendsUserName}}, bson.D{{Key: "friends_b_username", Value: ownerUserName}, {Key: "friend_a_username", Value: friendsUserName}}}}}

	err = friends_Repo.DB.Collection(entity.Friends).FindOne(context.TODO(), filter).Decode(friends)
	if err != nil {
		log.Println("error finding friends by his/her username")
		return nil, err
	}
	return friends, nil

}
