package user

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"go-friend-mgmt/cmd/internal/services/models"
	"testing"
	_ "github.com/lib/pq"
)

const (
	hostTest="localhost"
	portTest=5432
	userTest="postgres"
	passwordTest="123456789"
	dbnameTest="FriendManagement"
)

func ConnectionDBForTest() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		hostTest, portTest, userTest, passwordTest, dbnameTest)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("sssssss")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected !")
	return db
}

func TestCreateConnectionFriend(t *testing.T)  {
	db:=ConnectionDBForTest()
	defer db.Close()
	testCase:=[]struct{
		name string
		friends []string
		expectedResponse *models.Response
	}{
		{
			name: "success",
			friends: []string{"david@gmail.com","uresh@gmail.com"},
			expectedResponse: &models.Response{
				Success: true,
			},
		},
		{
			name: "email does not exists in db",
			friends: []string{"david@gmail.com","swap@gmail.com"},
			expectedResponse: &models.Response{
				Success: false,
			},
		},
		{
			name: "create failed because they are connected as friend",
			friends: []string{"david@gmail.com","uresh@gmail.com"},
			expectedResponse: &models.Response{
				Success: false,
			},
		},
		{
			name: "create failed because they are blocked",
			friends: []string{"uresh@gmail.com","kimsong@gmail.com"},
			expectedResponse: &models.Response{
				Success: false,
			},
		},
	}
	service:=ServiceImpl{DB: db}
	err:=insertConnectionFriend(service.DB)
	require.NoError(t, err)
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req:=models.FriendsList{
				Friends: tt.friends,
			}
			res, _:=service.CreateConnectionFriend(req)
			//require.NoError(t, err)
			require.Equal(t,tt.expectedResponse,res)
		})
	}
}

func TestReceiveFriendListByEmail(t *testing.T)  {
	db:=ConnectionDBForTest()
	defer db.Close()
	testCase:=[]struct{
		name string
		email string
		expectedResponse *models.ResponseFriend
	}{
		{
			name: "success",
			email: "david@gmail.com",
			expectedResponse: &models.ResponseFriend{
				Success: true,
				Friends: []string{"anna@gmail.com","swap@gmail.com","uresh@gmail.com"},
				Count: 3,
			},
		},
		{
			name: "invalid email",
			email: "davidgmail.com",
			expectedResponse: &models.ResponseFriend{
				Success: false,
				Count: 0,
			},
		},
		{
			name: "email does not exists in db",
			email: "lisa@gmail.com",
			expectedResponse: &models.ResponseFriend{
				Success: false,
				Count: 0,
			},
		},
	}
	service:=ServiceImpl{DB: db}
	err :=insertFriend(service.DB)
	require.NoError(t, err)
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req:=tt.email
			res, _:=service.ReceiveFriendListByEmail(req)
			require.NoError(t, err)
			require.Equal(t,tt.expectedResponse,res)
		})
	}
}

func TestReceiveCommonFriendList(t *testing.T)  {
	db:=ConnectionDBForTest()
	defer db.Close()
	testCase:=[]struct{
		name string
		friends models.Friend
		expectedResponse *models.ResponseFriend
	}{
		{
			name: "success",
			friends: models.Friend{
				UserEmail: "david@gmail.com",
				FriendEmail: "anna@gmail.com",
			},
			expectedResponse: &models.ResponseFriend{
				Success: true,
				Friends: []string{"swap@gmail.com","tom@gmail.com","uresh@gmail.com"},
				Count: 3,
			},
		},
		{
			name: "invalid email",
			friends: models.Friend{
				UserEmail: "davidgmail.com",
				FriendEmail: "anna@gmail.com",
			},
			expectedResponse: &models.ResponseFriend{
				Success: false,
				Count: 0,
			},
		},
		{
			name: "email does not exists in db",
			friends: models.Friend{
				UserEmail: "lisa@gmail.com",
				FriendEmail: "anna@gmail.com",
			},
			expectedResponse: &models.ResponseFriend{
				Success: false,
				Count: 0,
			},
		},
	}
	service:=ServiceImpl{DB: db}
	err :=insertCommonFriend(service.DB)
	require.NoError(t, err)
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req:=tt.friends
			res, _:=service.ReceiveCommonFriendList(req)
			require.NoError(t, err)
			require.Equal(t,tt.expectedResponse,res)
		})
	}
}

func TestSubscribeUpdateFromEmail(t *testing.T)  {
	db:=ConnectionDBForTest()
	defer db.Close()
	testCase:=[]struct{
		name string
		subscribe models.SubscribeUser
		expectedResponse *models.Response
	}{
		{
			name: "success with relationship not existed in db yet",
			subscribe: models.SubscribeUser{
				Requestor: "david@gmail.com",
				Target: "swap@gmail.com",
			},
			expectedResponse: &models.Response{
				Success: true,
			},
		},
		{
			name: "success with relationship already exist in db",
			subscribe: models.SubscribeUser{
				Requestor: "anna@gmail.com",
				Target: "david@gmail.com",
			},
			expectedResponse: &models.Response{
				Success: true,
			},
		},
		{
			name: "invalid email",
			subscribe: models.SubscribeUser{
				Requestor: "davidgmail.com",
				Target: "anna@gmail.com",
			},
			expectedResponse: &models.Response{
				Success: false,
			},
		},
		{
			name: "email does not exists in db",
			subscribe: models.SubscribeUser{
				Requestor: "lisa@gmail.com",
				Target: "anna@gmail.com",
			},
			expectedResponse: &models.Response{
				Success: false,
			},
		},
	}
	service:=ServiceImpl{DB: db}
	err :=insertSubscribeFriend(service.DB)
	require.NoError(t, err)
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req:=tt.subscribe
			res, _:=service.SubscribeUpdateFromEmail(req)
			require.NoError(t, err)
			require.Equal(t,tt.expectedResponse,res)
		})
	}
}

func TestBlockUpdateFromEmail(t *testing.T)  {
	db:=ConnectionDBForTest()
	defer db.Close()
	testCase:=[]struct{
		name string
		subscribe models.SubscribeUser
		expectedResponse *models.Response
	}{
		{
			name: "success with relationship not existed in db yet",
			subscribe: models.SubscribeUser{
				Requestor: "david@gmail.com",
				Target: "swap@gmail.com",
			},
			expectedResponse: &models.Response{
				Success: true,
			},
		},
		{
			name: "success with relationship already exist in db",
			subscribe: models.SubscribeUser{
				Requestor: "anna@gmail.com",
				Target: "david@gmail.com",
			},
			expectedResponse: &models.Response{
				Success: true,
			},
		},
		{
			name: "invalid email",
			subscribe: models.SubscribeUser{
				Requestor: "davidgmail.com",
				Target: "anna@gmail.com",
			},
			expectedResponse: &models.Response{
				Success: false,
			},
		},
		{
			name: "email does not exists in db",
			subscribe: models.SubscribeUser{
				Requestor: "lisa@gmail.com",
				Target: "anna@gmail.com",
			},
			expectedResponse: &models.Response{
				Success: false,
			},
		},
	}
	service:=ServiceImpl{DB: db}
	err :=insertBlockFriend(service.DB)
	require.NoError(t, err)
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req:=tt.subscribe
			res, _:=service.BlockUpdateFromEmail(req)
			require.NoError(t, err)
			require.Equal(t,tt.expectedResponse,res)
		})
	}
}

func TestGetAllSubscribeUpdateByEmail(t *testing.T)  {
	db:=ConnectionDBForTest()
	defer db.Close()
	testCase:=[]struct{
		name string
		retrieve models.RetrieveUpdate
		expectedResponse *models.SubscribeResponse
	}{
		{
			name: "success",
			retrieve: models.RetrieveUpdate{
				Sender: "david@gmail.com",
				Text: "Hello World! kate@gmail.com",
			},
			expectedResponse: &models.SubscribeResponse{
				Success: true,
				Recipients: []string{"anna@gmail.com","tom@gmail.com","uresh@gmail.com","kate@gmail.com"},
			},
		},
		{
			name: "invalid email",
			retrieve: models.RetrieveUpdate{
				Sender: "davidgmail.com",
				Text: "Hello World! kate@example",
			},
			expectedResponse: &models.SubscribeResponse{
				Success: false,
			},
		},
		{
			name: "email does not exists in db",
			retrieve: models.RetrieveUpdate{
				Sender: "lisa@gmail.com",
				Text: "Hello World! kate@example",
			},
			expectedResponse: &models.SubscribeResponse{
				Success: false,
			},
		},
	}
	service:=ServiceImpl{DB: db}
	err :=insertReceiveUpdate(service.DB)
	require.NoError(t, err)
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req:=tt.retrieve
			res, _:=service.GetAllSubscribeUpdateByEmail(req)
			require.NoError(t, err)
			require.Equal(t,tt.expectedResponse,res)
		})
	}
}

func insertConnectionFriend(db *sql.DB) error  {
	sqlStatement:=`truncate "User" cascade;
					truncate "UserRelationship" cascade;
					INSERT INTO "User" ("email")
 					VALUES ('david@gmail.com'),
 					('anna@gmail.com'),
 					('uresh@gmail.com');
					INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate")
 					VALUES ('anna@gmail.com','david@gmail.com',true,''),
 					('uresh@gmail.com','kimsong@gmail.com',false ,'BLOCK');`
	_,err:=db.Exec(sqlStatement)
	if err!=nil{
		fmt.Print(err)
		return err
	}
	return nil
}

func insertFriend(db *sql.DB) error  {
	sqlStatement:=`truncate "User" cascade;
					truncate "UserRelationship" cascade;
					INSERT INTO "User" ("email")
 					VALUES ('david@gmail.com');
					INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate")
 					VALUES ('anna@gmail.com','david@gmail.com',true,''),
 					('david@gmail.com','uresh@gmail.com',true ,''),
 					('david@gmail.com','swap@gmail.com',true ,'');`
	_,err:=db.Exec(sqlStatement)
	if err!=nil{
		return err
	}
	return nil
}

func insertCommonFriend(db *sql.DB) error  {
	sqlStatement:=`truncate "User" cascade;
					truncate "UserRelationship" cascade;
					INSERT INTO "User" ("email")
 					VALUES ('david@gmail.com'),
 							('anna@gmail.com');
					INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate")
 					VALUES ('anna@gmail.com','david@gmail.com',true,''),
 					('david@gmail.com','uresh@gmail.com',true ,''),
 					('tom@gmail.com','david@gmail.com',true ,''),
 					('david@gmail.com','swap@gmail.com',true ,''),
 					('anna@gmail.com','jerry@gmail.com',true,''),
 					('swap@gmail.com','anna@gmail.com',true ,''),
 					('tom@gmail.com','anna@gmail.com',true ,''),
 					('uresh@gmail.com','anna@gmail.com',true ,'');`
	_,err:=db.Exec(sqlStatement)
	if err!=nil{
		return err
	}
	return nil
}

func insertSubscribeFriend(db *sql.DB) error  {
	sqlStatement:=`truncate "User" cascade;
					truncate "UserRelationship" cascade;
					INSERT INTO "User" ("email")
 					VALUES ('david@gmail.com'),
 							('swap@gmail.com'),
 							('anna@gmail.com');
					INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate")
 					VALUES ('anna@gmail.com','david@gmail.com',true,'');`
	_,err:=db.Exec(sqlStatement)
	if err!=nil{
		return err
	}
	return nil
}

func insertBlockFriend(db *sql.DB) error  {
	sqlStatement:=`truncate "User" cascade;
					truncate "UserRelationship" cascade;
					INSERT INTO "User" ("email")
 					VALUES ('david@gmail.com'),
 							('swap@gmail.com'),
 							('anna@gmail.com');
					INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate")
 					VALUES ('anna@gmail.com','david@gmail.com',true,'');`
	_,err:=db.Exec(sqlStatement)
	if err!=nil{
		return err
	}
	return nil
}

func insertReceiveUpdate(db *sql.DB) error  {
	sqlStatement:=`truncate "User" cascade;
					truncate "UserRelationship" cascade;
					INSERT INTO "User" ("email")
 					VALUES ('david@gmail.com'),
 							('kate@gmail.com');
					INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate")
 					VALUES ('anna@gmail.com','david@gmail.com',true,''),
 					('david@gmail.com','uresh@gmail.com',true ,'SUBSCRIBE'),
 					('tom@gmail.com','david@gmail.com',false ,'SUBSCRIBE'),
 					('david@gmail.com','swap@gmail.com',false ,'BLOCK');`
	_,err:=db.Exec(sqlStatement)
	if err!=nil{
		return err
	}
	return nil
}