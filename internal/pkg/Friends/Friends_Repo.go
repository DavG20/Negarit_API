package friends

import (
	"context"
	"errors"
	"fmt"
	"log"

	message "github.com/DavG20/Negarit_API/internal/pkg/Message"
	usermodel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
	userservice "github.com/DavG20/Negarit_API/internal/pkg/User/User_Service"
	"github.com/DavG20/Negarit_API/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IFriendsRepo interface {
	CreateFriendship(message.Message) (*Friends, error)
	AreTheyFriend(userName, friendsUserName string) bool
	GetFriendByUserName(ownerUserName, friendsUserName string) (*Friends, error)
	GetAllFriends(finderUserName string) ([]*usermodel.DBResponse, error)
	AppendMessage(*Friends) (*Friends, error)
	DeleteFriends(userName, friendsUserName string) bool
	BlockFriend(string, *Friends) bool
}

type Friends_Repo struct {
	DB *mongo.Database
}

func NewFriendsRepo(db *mongo.Database) *Friends_Repo {
	return &Friends_Repo{
		DB: db,
	}
}

// create friendship function
//
//	create friendship from the message
//
// for message the input is reciever's id and message's content
func (friends_Repo *Friends_Repo) CreateFriendship(message message.Message) (friends *Friends, err error) {

	friends = &Friends{}
	friends_Id := ""
	friends.Friend_A_UserName = message.Sender_UserName
	friends.Friend_B_UserName = message.Receiver_UserName
	friends.Message = append(friends.Message, message)
	friends.Block_By_A = false
	friends.Block_By_B = false

	res, err := friends_Repo.DB.Collection(entity.Friends).InsertOne(context.TODO(), friends)
	if err != nil {
		return nil, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		friends_Id = entity.GetIdFromInsertedObjectId(oid)
	}
	friends.Friend_Id = friends_Id
	return friends, nil

}

func (friends_Repo *Friends_Repo) AppendMessage(friends *Friends) (*Friends, error) {

	filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "friend_a_username", Value: friends.Friend_A_UserName}, {Key: "friend_b_username", Value: friends.Friend_B_UserName}}, bson.D{{Key: "friend_b_username", Value: friends.Friend_A_UserName}, {Key: "friend_a_username", Value: friends.Friend_B_UserName}}}}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "message", Value: friends.Message}}}}
	_, err := friends_Repo.DB.Collection(entity.Friends).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return friends, nil

}

// check if they are friends when the  user sends message in the first time
// if they are friend return true if not return false

func (friends_Repo *Friends_Repo) AreTheyFriend(userName, friendsUserName string) bool {
	friends := &Friends{}

	filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "friend_a_username", Value: userName}, {Key: "friend_b_username", Value: friendsUserName}}, bson.D{{Key: "friend_b_username", Value: userName}, {Key: "friend_a_username", Value: friendsUserName}}}}}
	err := friends_Repo.DB.Collection(entity.Friends).FindOne(context.TODO(), filter).Decode(&friends)
	if err != nil {
		log.Println("error finding friends")
		return false
	}
	return true
}

// to get friends by his/her username  and return friends
// they must be friends already

func (friends_Repo *Friends_Repo) GetFriendByUserName(ownerUserName, friendsUserName string) (*Friends, error) {
	friends := &Friends{}
	// filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "friend_b_username", Value: ownerUserName}, {Key: "friends_b_username", Value: friendsUserName}}, bson.D{{Key: "friends_a_username", Value: friendsUserName}, {Key: "friend_b_username", Value: ownerUserName}}}}}
	filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "friend_a_username", Value: ownerUserName}, {Key: "friend_b_username", Value: friendsUserName}}, bson.D{{Key: "friend_a_username", Value: friendsUserName}, {Key: "friend_b_username", Value: ownerUserName}}}}}
	err := friends_Repo.DB.Collection(entity.Friends).FindOne(context.TODO(), filter).Decode(friends)

	if err != nil {
		log.Println("error finding friends by his/her username", friends)
		return nil, err
	}
	// fmt.Println(friends, "friend in repo,get friends by username line 148")
	return friends, nil

}

func (friends_Repo *Friends_Repo) GetAllFriends(finderUserName string) (Users []*usermodel.DBResponse, err error) {
	// friends := []Friends{}
	userServ := userservice.UserService{}
	filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "friend_a_username", Value: finderUserName}}, bson.D{{Key: "friend_b_username", Value: finderUserName}}}}}

	cursor, err := friends_Repo.DB.Collection(entity.Friends).Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("error in here getAllfriends line 130")
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		friend := &Friends{}
		friends_userName := ""
		cursor.Decode(friend)
		if friend.Friend_Id == "" {
			continue
		}
		if friend.Friend_A_UserName == finderUserName {
			friends_userName = friend.Friend_B_UserName
		} else {
			friends_userName = friend.Friend_A_UserName
		}

		// user:=user_service.GetUserByUserName(friends_userName)
		user := userServ.GetUserByUserName(friends_userName)
		if user == nil {
			continue
		}
		dbUser := userServ.GetSecuredUser(user)
		Users = append(Users, dbUser)

	}
	if len(Users) == 0 {
		return nil, errors.New("u haven't any friend yet")
	}
	return Users, nil

}

func (friends_Repo *Friends_Repo) GetUser(username string) (user *usermodel.User, err error) {
	filter := bson.D{{Key: "username", Value: username}}

	err = friends_Repo.DB.Collection(entity.User).FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (friends_Repo *Friends_Repo) CheckUserNameExist(userName string) bool {
	_, err := friends_Repo.GetUser(userName)
	if err != nil {
		return false
	}
	return true

}
func (friend_Repo *Friends_Repo) UpdateFriendsID(friend *Friends) bool {
	filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "friend_a_username", Value: friend.Friend_A_UserName}, {Key: "friends_b_username", Value: friend.Friend_B_UserName}}, bson.D{{Key: "friends_a_username", Value: friend.Friend_B_UserName}, {Key: "friend_b_username", Value: friend.Friend_A_UserName}}}}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "id", Value: friend.Friend_Id}}}}
	fmt.Println(friend.Friend_Id, "new id")
	_, err := friend_Repo.DB.Collection(entity.Friends).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("error updating friends id ")
		return false
	}
	return true
}

func (friends_Repo *Friends_Repo) DeleteFriends(userName, friendsUserName string) bool {
	filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "friend_a_username", Value: userName}, {Key: "friend_b_username", Value: friendsUserName}}, bson.D{{Key: "friend_a_username", Value: friendsUserName}, {Key: "friend_b_username", Value: userName}}}}}

	_, err := friends_Repo.DB.Collection(entity.Friends).DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println("error deleting friend")
		return false
	}

	return true
}

func (friends_Repo *Friends_Repo) BlockFriend(userName string, friends *Friends) bool {
	if friends == nil {
		fmt.Println("invalid friend")
		return false
	}
	filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "friend_a_username", Value: friends.Friend_A_UserName}, {Key: "friend_b_username", Value: friends.Friend_B_UserName}}, bson.D{{Key: "friend_a_username", Value: friends.Friend_B_UserName}, {Key: "friend_b_username", Value: friends.Friend_A_UserName}}}}}
	// fr := &Friends{}
	if userName == friends.Friend_A_UserName {
		fmt.Println("in")
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "blocka", Value: !friends.Block_By_A}}}}
		_, err := friends_Repo.DB.Collection(entity.Friends).UpdateOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println("error finding friend")
			return false
		}
		return true
	} else {
		fmt.Println("in else")
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "blockb", Value: true}}}}
		_, err := friends_Repo.DB.Collection(entity.Friends).ReplaceOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println("error finding friend")
			return false
		}
		return true
	}
}
