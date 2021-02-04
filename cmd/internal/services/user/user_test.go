package user

import (
	"database/sql"
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
		{
			name: "email already exists in db",
			email: "john@gmail.com",
			expectedResponse: &models.Response{
				Success: false,
			},
		},
	}
	service:=ServiceImpl{DB: db}
	err:=insertUser(service.DB)
	require.NoError(t,err)
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

func insertUser(db *sql.DB) error  {
	sqlStatement:=`
					INSERT INTO "User" ("email")
 					VALUES ('john@gmail.com');
					`
	_,err:=db.Exec(sqlStatement)
	if err!=nil{
		return err
	}
	return nil
}

