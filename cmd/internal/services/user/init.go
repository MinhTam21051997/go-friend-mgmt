package user

import (
	"database/sql"
	"go-friend-mgmt/cmd/internal/services/models"
)

type Service interface {
	RetrieveByID(ID int) (*models.Users, error)
	CreateConnectionFriend(friend models.Friend) (int64,*models.Response )
	ReceiveFriendListByEmail(email string) (*models.ResponseFriend,error)
	ReceiveCommonFriendList(friend models.Friend) (*models.ResponseFriend,error)
	SubscribeUpdateFromEmail(subscribeUser models.SubscribeUser) (*models.Response,error)
	BlockUpdateFromEmail(subscribeUser models.SubscribeUser) (*models.Response,error)
	GetAllSubscribeUpdateByEmail(retrieve models.RetrieveUpdate) (*models.SubscribeResponse,error)
}

type ServiceImpl struct {
	DB *sql.DB
}

type ManagerDB struct {
	service ServiceImpl
}