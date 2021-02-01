package router

import (
	"github.com/go-chi/chi"
	api "go-friend-mgmt/cmd/api/router/rest_api"
	"go-friend-mgmt/cmd/internal/services/user"
	"net/http"
)

type RouterHandler struct {
	ProductService user.Service
}

func (handler RouterHandler)InitializeRoutes() http.Handler{
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Post("/user/makefriend",api.CreateConnectionFriend(handler.ProductService))
		r.Post("/user/createuser",api.CreateUser(handler.ProductService))
		r.Post("/user/friends",api.ReceiveFriendListByEmail(handler.ProductService))
		r.Post("/user/commonfriends",api.ReceiveCommonFriendList(handler.ProductService))
		r.Post("/user/subscribe",api.SubscribeUpdateFromEmail(handler.ProductService))
		r.Post("/user/block",api.BlockUpdateFromEmail(handler.ProductService))
		r.Post("/user/emailssubscribe",api.GetAllSubscribeUpdateByEmail(handler.ProductService))
	})
	return r
}