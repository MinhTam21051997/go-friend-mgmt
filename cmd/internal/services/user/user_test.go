package user

import (
	"testing"
	"github.com/stretchr/testify/require"
	"go-friend-mgmt/cmd/internal/services/models"
)

func TestCreateUser(t *testing.T)  {
	db:=ConnectionDBForTest()
	defer db.Close()
	testCase:=[]struct{
		name string
		email string
		expectedResponse *models.Response
	}{
		{
			name: "success",
			email: "david@gmail.com",
			expectedResponse: &models.Response{
				Success: true,
			},
		},
		{
			name: "invalid email",
			email: "davidgmail.com",
			expectedResponse: &models.Response{
				Success: false,
			},
		},
	}
	service:=ServiceImpl{DB: db}
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req:=models.EmailUser{
				Email: tt.email,
			}
			res, _:=service.CreateUser(req)
			//require.NoError(t, err)
			require.Equal(t,tt.expectedResponse,res)
		})
	}
}

