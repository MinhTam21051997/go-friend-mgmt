package rest_api

import (
	"encoding/json"
	"go-friend-mgmt/cmd/internal/services/models"
	"go-friend-mgmt/cmd/internal/services/user"
	"net/http"
)

func CreateUser(service user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var emailUser models.EmailUser
		err:=json.NewDecoder(r.Body).Decode(&emailUser)
		if err!=nil{
			responseWithJson(w, http.StatusBadRequest, RenderBadRequest(err))
			return
		}
		req:=models.EmailUser{
			Email: emailUser.Email,
		}
		response, err := service.CreateUser(req)
		if err!=nil{
			responseWithJson(w,http.StatusInternalServerError,ServerErrorRender(err))
			return
		}
		responseWithJson(w, http.StatusOK, response)
	}
}
