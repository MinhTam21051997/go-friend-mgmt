package main

import (
	"fmt"
	"go-friend-mgmt/cmd/api/router"
	"go-friend-mgmt/cmd/internal/services/database"
	"go-friend-mgmt/cmd/internal/services/user"
	"net/http"
)

func main(){
	db:=database.ConnectionDB()
	defer db.Close()
	router := router.RouterHandler{
		ProductService: user.ServiceImpl{
			DB: db,
		},
	}

	err := http.ListenAndServe(":3060", router.InitializeRoutes())
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}