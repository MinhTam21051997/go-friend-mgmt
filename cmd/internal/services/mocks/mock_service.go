package mocks

import (
	"go-friend-mgmt/cmd/internal/services/models"
	"github.com/stretchr/testify/mock"
)
type ServiceMock struct {
	mock.Mock
}

func (s ServiceMock) CreateConnectionFriend(req models.FriendsList) (*models.Response,error) {
	returnVals := s.Called(req)
	r0 := returnVals.Get(0).(models.Response)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return &r0, r1
}

func (s ServiceMock) CreateUser(req models.EmailUser) (*models.Response,error) {
	returnVals := s.Called(req)
	r0 := returnVals.Get(0).(models.Response)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return &r0, r1
}


func (s ServiceMock) ReceiveFriendListByEmail(req string) (*models.ResponseFriend,error) {
	returnVals := s.Called(req)
	r0 := returnVals.Get(0).(models.ResponseFriend)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return &r0, r1
}

func (s ServiceMock) ReceiveCommonFriendList(req models.Friend) (*models.ResponseFriend,error) {
	returnVals := s.Called(req)
	r0 := returnVals.Get(0).(models.ResponseFriend)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return &r0, r1
}

func (s ServiceMock) SubscribeUpdateFromEmail(req models.SubscribeUser) (*models.Response,error) {
	returnVals := s.Called(req)
	r0 := returnVals.Get(0).(models.Response)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return &r0, r1
}

func (s ServiceMock) BlockUpdateFromEmail(req models.SubscribeUser) (*models.Response,error) {
	returnVals := s.Called(req)
	r0 := returnVals.Get(0).(models.Response)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return &r0, r1
}

func (s ServiceMock) GetAllSubscribeUpdateByEmail(req models.RetrieveUpdate) (*models.SubscribeResponse,error) {
	returnVals := s.Called(req)
	r0 := returnVals.Get(0).(models.SubscribeResponse)
	var r1 error
	if returnVals.Get(1) != nil {
		r1 = returnVals.Get(1).(error)
	}
	return &r0, r1
}


