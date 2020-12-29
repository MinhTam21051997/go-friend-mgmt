package rest_api

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"go-friend-mgmt/cmd/internal/services/models"
	"go-friend-mgmt/cmd/internal/services/user"
	"log"
	"net/http"
	"strconv"
)

func GetUserById(service user.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		strID := chi.URLParam(r, "id")
		log.Println("aaaaa")
		id, err := strconv.Atoi(strID)
		if err == nil {
			// handle err
		}

		user, err := service.RetrieveByID(id)
		if err != nil {
			// handle err
		}

		responseWithJson(w, http.StatusOK, user)
	}
}

func CreateConnectionFriend(service user.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var friendsList models.FriendsList
		err:=json.NewDecoder(r.Body).Decode(&friendsList)
		if err!=nil{
			responseWithJson(w, http.StatusInternalServerError, err.Error())
			return
		}
		log.Println("friendlist", friendsList)
		//sliceFriend:=strings.Split(friend,",")
		if len(friendsList.Friends) != 2{
			responseWithJson(w, http.StatusInternalServerError, errors.New("Only receive two email !"))
		}

		friendModel:=models.Friend{
			UserEmail: friendsList.Friends[0],
			FriendEmail: friendsList.Friends[1],
		}
		id,response := service.CreateConnectionFriend(friendModel)
		log.Println("id ", id)
		if id < 0 {
			responseWithJson(w, http.StatusInternalServerError, response)
		}

		responseWithJson(w, http.StatusOK, response)
	}
}

func ReceiveFriendListByEmail(service user.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var emailUser models.EmailUser
		json.NewDecoder(r.Body).Decode(&emailUser)
		resFriend, err := service.ReceiveFriendListByEmail(emailUser.Email)
		if err!=nil{
			responseWithJson(w, http.StatusInternalServerError, err.Error())
			return
		}
		responseWithJson(w, http.StatusOK, resFriend)
	}
}

func ReceiveCommonFriendList(service user.Service) func(w http.ResponseWriter, r *http.Request)  {
	return func(w http.ResponseWriter, r *http.Request) {
		var friendsList models.FriendsList
		err:=json.NewDecoder(r.Body).Decode(&friendsList)
		if err!=nil{
			responseWithJson(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(friendsList.Friends) != 2{
			responseWithJson(w, http.StatusInternalServerError, &models.ResponseFriend{
				Success: false,
				Friends: nil,
				Count: 0,
				Message: "Need to input two email !",
			})
			return
		}

		friendModel:=models.Friend{
			UserEmail: friendsList.Friends[0],
			FriendEmail: friendsList.Friends[1],
		}
		response, err := service.ReceiveCommonFriendList(friendModel)
		if err !=nil {
			responseWithJson(w, http.StatusInternalServerError, err.Error())
			return
		}

		responseWithJson(w, http.StatusOK, response)
	}
}

func SubscribeUpdateFromEmail(service user.Service) func(w http.ResponseWriter, r *http.Request)  {
	return func(w http.ResponseWriter, r *http.Request) {
		var subscribeUser models.SubscribeUser
		err:=json.NewDecoder(r.Body).Decode(&subscribeUser)
		if err!=nil{
			responseWithJson(w, http.StatusInternalServerError, err.Error())
			return
		}
		response, err := service.SubscribeUpdateFromEmail(subscribeUser)
		if err !=nil {
			responseWithJson(w, http.StatusInternalServerError, err.Error())
			return
		}

		responseWithJson(w, http.StatusOK, response)
	}
}

func BlockUpdateFromEmail(service user.Service) func(w http.ResponseWriter, r *http.Request)  {
	return func(w http.ResponseWriter, r *http.Request) {
		var subscribeUser models.SubscribeUser
		err:=json.NewDecoder(r.Body).Decode(&subscribeUser)
		if err!=nil{
			responseWithJson(w, http.StatusInternalServerError, err.Error())
			return
		}
		response, err := service.BlockUpdateFromEmail(subscribeUser)
		if err !=nil {
			responseWithJson(w, http.StatusInternalServerError, err.Error())
			return
		}

		responseWithJson(w, http.StatusOK, response)
	}
}

func GetAllSubscribeUpdateByEmail(service user.Service) func(w http.ResponseWriter, r *http.Request)  {
	return func(w http.ResponseWriter, r *http.Request) {
		var retrieveUpdate models.RetrieveUpdate
		err:=json.NewDecoder(r.Body).Decode(&retrieveUpdate)
		log.Println("retrieveUpdate :", retrieveUpdate)
		if err!=nil{
			responseWithJson(w, http.StatusInternalServerError, err.Error())
			return
		}
		response, err := service.GetAllSubscribeUpdateByEmail(retrieveUpdate)
		if err !=nil {
			responseWithJson(w, http.StatusInternalServerError, err.Error())
			return
		}

		responseWithJson(w, http.StatusOK, response)
	}
}