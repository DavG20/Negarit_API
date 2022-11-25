package friends

import (
	message "github.com/DavG20/Negarit_API/internal/pkg/Message"
	usermodel "github.com/DavG20/Negarit_API/internal/pkg/User/User_Model"
)

type IFriendsService interface {
	CreateFriendship(message.Message) (*Friends, error)
	AreTheyFriend(userName, friendsUserName string) bool
	GetFriendByUserName(finderUserName, friendsUserName string) (*Friends, error)
	GetAllFriends(finderUserName string) ([]*usermodel.DBResponse, error)
	AppendMessage(*Friends) bool
	DeleteFriends(userName, friendsUserName string) bool
	BlockFriend(userName, friendsUserNAme string) (*Friends, error)
}

type FriendsService struct {
	IFriendsRepo IFriendsRepo
}

func NewFriendsService(iFriendsRepo IFriendsRepo) FriendsService {
	return FriendsService{
		IFriendsRepo: iFriendsRepo,
	}
}

func (friendsService *FriendsService) CreateFriendship(message message.Message) (*Friends, error) {
	return friendsService.IFriendsRepo.CreateFriendship(message)
}

func (friendsService *FriendsService) AreTheyFriend(userName, friendsUserName string) bool {
	return friendsService.IFriendsRepo.AreTheyFriend(userName, friendsUserName)
}

func (friendsService *FriendsService) GetFriendByUserName(finderUserName, friendsUserName string) (*Friends, error) {
	return friendsService.IFriendsRepo.GetFriendByUserName(finderUserName, friendsUserName)
}

func (friendsService *FriendsService) GetAllFriends(finderUserName string) (*usermodel.DBResponse, error) {
	return friendsService.GetAllFriends(finderUserName)
}

func (friendsService *FriendsService) AppendMessage(friends *Friends) bool {
	_, err := friendsService.IFriendsRepo.AppendMessage(friends)
	if err != nil {
		return false
	}
	return true
}

func (friendsService *FriendsService) DeleteFriends(userName, friendsUserName string) bool {
	return friendsService.IFriendsRepo.DeleteFriends(userName, friendsUserName)
}

func (friendsService *FriendsService) BlockFriend(userName string,friends *Friends) ( bool) {
	return friendsService.IFriendsRepo.BlockFriend(userName, friends)
	
}
