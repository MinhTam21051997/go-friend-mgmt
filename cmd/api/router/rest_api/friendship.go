package rest_api

import (
	"encoding/json"
	"go-friend-mgmt/cmd/internal/services/models"
	"go-friend-mgmt/cmd/internal/services/user"
	"net/http"
)

func CreateConnectionFriend(service user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var friendsList models.FriendsList
		err:=json.NewDecoder(r.Body).Decode(&friendsList)
		if err!=nil{
			responseWithJson(w, http.StatusBadRequest, RenderBadRequest(err))
			return
		}
		response, err := service.CreateConnectionFriend(friendsList)
		if err!=nil{
			responseWithJson(w,http.StatusInternalServerError,ServerErrorRender(err))
			return
		}
		responseWithJson(w, http.StatusOK, response)
	}
}

func ReceiveFriendListByEmail(service user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var emailUser models.EmailUser
		err:= json.NewDecoder(r.Body).Decode(&emailUser)
		if err!=nil{
			responseWithJson(w, http.StatusBadRequest, RenderBadRequest(err))
			return
		}
		resFriend, err := service.ReceiveFriendListByEmail(emailUser.Email)
		if err!=nil{
			responseWithJson(w, http.StatusInternalServerError, ServerErrorRender(err))
			return
		}
		responseWithJson(w, http.StatusOK, resFriend)
	}
}

func ReceiveCommonFriendList(service user.Service) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		var friendsList models.FriendsList
		err:=json.NewDecoder(r.Body).Decode(&friendsList)
		if err!=nil{
			responseWithJson(w, http.StatusBadRequest, RenderBadRequest(err))
			return
		}
		friendModel:=models.Friend{
			UserEmail: friendsList.Friends[0],
			FriendEmail: friendsList.Friends[1],
		}
		response, err := service.ReceiveCommonFriendList(friendModel)
		if err !=nil {
			responseWithJson(w, http.StatusInternalServerError, ServerErrorRender(err))
			return
		}
		responseWithJson(w, http.StatusOK, response)
	}
}

func SubscribeUpdateFromEmail(service user.Service) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		var subscribeUser models.SubscribeUser
		err:=json.NewDecoder(r.Body).Decode(&subscribeUser)
		if err!=nil{
			responseWithJson(w, http.StatusBadRequest, RenderBadRequest(err))
			return
		}
		response, err := service.SubscribeUpdateFromEmail(subscribeUser)
		if err !=nil {
			responseWithJson(w, http.StatusInternalServerError, ServerErrorRender(err))
			return
		}
		responseWithJson(w, http.StatusOK, response)
	}
}

func BlockUpdateFromEmail(service user.Service) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		var subscribeUser models.SubscribeUser
		err:=json.NewDecoder(r.Body).Decode(&subscribeUser)
		if err!=nil{
			responseWithJson(w, http.StatusBadRequest, RenderBadRequest(err))
			return
		}
		response, err := service.BlockUpdateFromEmail(subscribeUser)
		if err !=nil {
			responseWithJson(w, http.StatusInternalServerError, ServerErrorRender(err))
			return
		}
		responseWithJson(w, http.StatusOK, response)
	}
}

func GetAllSubscribeUpdateByEmail(service user.Service) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		var retrieveUpdate models.RetrieveUpdate
		err:=json.NewDecoder(r.Body).Decode(&retrieveUpdate)
		if err!=nil{
			responseWithJson(w, http.StatusBadRequest, RenderBadRequest(err))
			return
		}
		response, err := service.GetAllSubscribeUpdateByEmail(retrieveUpdate)
		if err !=nil {
			responseWithJson(w, http.StatusInternalServerError, ServerErrorRender(err))
			return
		}
		responseWithJson(w, http.StatusOK, response)
	}
}