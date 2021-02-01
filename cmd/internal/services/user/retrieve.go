package user

import (
	"database/sql"
	"errors"
	"fmt"
	"go-friend-mgmt/cmd/internal/services/models"
	"go-friend-mgmt/cmd/internal/services/utils"
	"log"
	"regexp"
)

func (s ServiceImpl) CreateConnectionFriend(friendList models.FriendsList) (*models.Response,error ) {
	response:=&models.Response{}
	if len(friendList.Friends) !=2{
		return response, errors.New("The request need to input two email")
	}
	friend:=models.Friend{
		UserEmail: friendList.Friends[0],
		FriendEmail: friendList.Friends[1],
	}
	err:=utils.VaditationEmail(friend.UserEmail)
	if err!=nil{
		return response, err
	}
	err=utils.VaditationEmail(friend.FriendEmail)
	if err!=nil{
		return response, err
	}
	err=CheckUserExistsInDb(s.DB,friend.UserEmail)
	if err!=nil{
		return response,err
	}
	err=CheckUserExistsInDb(s.DB,friend.FriendEmail)
	if err!=nil{
		return response,err
	}
	err=CheckIsFriendOrBlockInDb(s.DB,friend.UserEmail,friend.FriendEmail)
	if err!=nil{
		return response,err
	}else {
		err:=CreateConnection(s.DB,friend.UserEmail, friend.FriendEmail)
		if err!=nil{
			response.Success=false
			return response, err
		}
	}
	response.Success=true
	return response,nil
}

func (s ServiceImpl)CreateUser(email models.EmailUser) (*models.Response, error)  {
	response:=&models.Response{}
	err:=utils.VaditationEmail(email.Email)
	if err!=nil{
		return response,err
	}
	err=CreateUserByEmail(s.DB, email.Email)
	if err!=nil{
		return response, err
	}
	response.Success=true
	return response,nil
}


func (s ServiceImpl) ReceiveFriendListByEmail(email string) (*models.ResponseFriend, error)  {
	responseFriend:= &models.ResponseFriend{}
	err:=utils.VaditationEmail(email)
	if err!=nil{
		return responseFriend,err
	}
	err=CheckUserExistsInDb(s.DB,email)
	if err!=nil{
		return responseFriend, err
	}
	res, err:=GetFriendListByEmail(s.DB, email)
	if err!=nil{
		return responseFriend, nil
	}
	responseFriend.Success=true
	responseFriend.Friends=res.Friends
	responseFriend.Count=int64(len(res.Friends))
	return responseFriend, nil
}

func (s ServiceImpl) ReceiveCommonFriendList(friend models.Friend) (*models.ResponseFriend,error)  {
	responseFriend:= &models.ResponseFriend{}
	err:=utils.VaditationEmail(friend.UserEmail)
	if err!=nil{
		return responseFriend,err
	}
	err=utils.VaditationEmail(friend.FriendEmail)
	if err!=nil{
		return responseFriend,err
	}
	err=CheckUserExistsInDb(s.DB,friend.UserEmail)
	if err!=nil{
		return responseFriend,err
	}
	err=CheckUserExistsInDb(s.DB,friend.FriendEmail)
	if err!=nil{
		return responseFriend,err
	}
	res, err:=GetCommonFriend(s.DB, friend.UserEmail,friend.FriendEmail)
	if err!=nil{
		responseFriend.Success=false
		return responseFriend,nil
	}
	responseFriend.Success=true
	responseFriend.Friends=res.Friends
	responseFriend.Count=int64(len(res.Friends))
	return responseFriend, nil
}

func (s ServiceImpl) SubscribeUpdateFromEmail(subscribeUser models.SubscribeUser) (*models.Response,error)  {
	response:= &models.Response{}
	if subscribeUser.Requestor==subscribeUser.Target {
		return response, errors.New("The two emails cannot be the same!")
	}
	err:=utils.VaditationEmail(subscribeUser.Requestor)
	if err!=nil{
		return response,err
	}
	err=utils.VaditationEmail(subscribeUser.Target)
	if err!=nil{
		return response,err
	}
	err=CheckUserExistsInDb(s.DB,subscribeUser.Requestor)
	if err!=nil{
		return response,err
	}
	err=CheckUserExistsInDb(s.DB,subscribeUser.Target)
	if err!=nil{
		return response,err
	}

	res, err:=CheckRelationshipExistsInDb(s.DB,subscribeUser.Requestor, subscribeUser.Target)
	if err!=nil{
		log.Println("Error ", err.Error())
		return response,err
	}
	if res < 1{
		err:=InsertRelationship(s.DB, subscribeUser.Requestor, subscribeUser.Target,"SUBSCRIBE")
		if  err!=nil{
			return response, err
		}
		response.Success=true
		return response, nil
	}
	err=UpdateRelationship(s.DB, subscribeUser.Requestor, subscribeUser.Target,"SUBSCRIBE")
	if  err!=nil{
		return response, err
	}
	response.Success=true
	return response, nil
}

func (s ServiceImpl) BlockUpdateFromEmail(subscribeUser models.SubscribeUser) (*models.Response,error)  {
	response:= &models.Response{}
	if subscribeUser.Requestor==subscribeUser.Target {
		return response, errors.New("The two emails cannot be the same!")
	}
	err:=utils.VaditationEmail(subscribeUser.Requestor)
	if err!=nil{
		return response,err
	}
	err=utils.VaditationEmail(subscribeUser.Target)
	if err!=nil{
		return response,err
	}
	err=CheckUserExistsInDb(s.DB,subscribeUser.Requestor)
	if err!=nil{
		return response,err
	}
	err=CheckUserExistsInDb(s.DB,subscribeUser.Target)
	if err!=nil{
		return response,err
	}
	res, err:=CheckRelationshipExistsInDb(s.DB,subscribeUser.Requestor, subscribeUser.Target)
	if err!=nil{
		log.Println("Error ", err.Error())
		return response, nil
	}
	if res < 1{
		err:=InsertRelationship(s.DB, subscribeUser.Requestor, subscribeUser.Target,"BLOCK")
		if  err!=nil{
			return response, err
		}
		response.Success=true
		return response, nil
	}
	err=UpdateRelationship(s.DB, subscribeUser.Requestor, subscribeUser.Target,"BLOCK")
	if  err!=nil{
		return response, err
	}
	response.Success=true
	return response, nil
}

func (s ServiceImpl) GetAllSubscribeUpdateByEmail(retrieve models.RetrieveUpdate) (*models.SubscribeResponse,error)  {
	subscribeResponse:= &models.SubscribeResponse{}
	err:=utils.VaditationEmail(retrieve.Sender)
	if err!=nil{
		return subscribeResponse,err
	}
	err=CheckUserExistsInDb(s.DB,retrieve.Sender)
	if err!=nil{
		return subscribeResponse,err
	}
	res, err:=GetAllSubscriberByEmail(s.DB, retrieve.Sender, retrieve.Text)
	if err!=nil{
		return subscribeResponse,err
	}
	subscribeResponse.Success=true
	subscribeResponse.Recipients=res.Friends
	return subscribeResponse, nil
}

func CheckUserExistsInDb(db *sql.DB, email string) error {
	var count int64
	err:=utils.VaditationEmail(email)
	if err!=nil{
		return err
	}else {
		err = db.QueryRow(`SELECT COUNT(*) FROM "User" WHERE "email" = $1`, email).Scan(&count)
		if err!=nil{
			return err
		}else {
			if count!=1{
				return errors.New("Email does not exists")
			}
		}
	}
	return nil
}

func CheckRelationshipExistsInDb(db *sql.DB, userEmail string, friendEmail string) (int,error) {
	var count int
	err:=utils.VaditationEmail(userEmail)
	if err!=nil{
		return count,err
	}
	err=utils.VaditationEmail(friendEmail)
	if err!=nil{
		return count,err
	}
	err=CheckUserExistsInDb(db,userEmail)
	if err!=nil{
		return count,err
	}
	err=CheckUserExistsInDb(db,friendEmail)
	if err!=nil{
		return count,err
	}
	err = db.QueryRow(`SELECT COUNT(*) FROM "UserRelationship" WHERE "UserEmail" = $1 and "FriendEmail" = $2`, userEmail, friendEmail).Scan(&count)
	if err!=nil{
		return count,err
	}
	return count,nil
}

func CheckIsFriendOrBlockInDb(db *sql.DB, userEmail string, friendEmail string) error {
	err:=utils.VaditationEmail(userEmail)
	if err!=nil{
		return err
	}
	err=utils.VaditationEmail(friendEmail)
	if err!=nil{
		return err
	}
	sqlQuery:=`SELECT * FROM "UserRelationship" WHERE ("UserEmail" = $1 and "FriendEmail" = $2) or ("FriendEmail" = $1 and "UserEmail" = $2)`
	rows,err := db.Query(sqlQuery, userEmail, friendEmail)
	if err!=nil{
		log.Println("[CheckIsFriendOrBlockInDb] query error : ", err.Error())
		return err
	}
	defer rows.Close()
	var relationships []models.UserRelationship
	for rows.Next(){
		var rel models.UserRelationship
		err := rows.Scan(&rel.RelationshipId,&rel.UserEmail,&rel.FriendEmail,&rel.IsFriend,&rel.StatusUpdate)
		if err!=nil{
			log.Fatalf("[CheckIsFriendOrBlock] Cannot scan to UserRelationship model %v ", err.Error() )
			return err
		}
		relationships=append(relationships,rel)
	}
	for i:=0;i<len(relationships);i++{
		if relationships[i].StatusUpdate=="BLOCK" || relationships[i].IsFriend==true{
			return errors.New("They were blocked or had friends")
		}
	}
	return nil
}

func CreateUserByEmail(db *sql.DB, emailUser string) error {
	sqlStatement:=`INSERT INTO "User" ("email") VALUES ($1)`
	_,err:=db.Exec(sqlStatement,emailUser)
	if err!=nil{
		return err
	}
	return nil
}

func CreateConnection(db *sql.DB, emailUser string, friendEmail string) error{
	var id int
	userRelationship:= models.UserRelationship{
		UserEmail: emailUser,
		FriendEmail: friendEmail,
		IsFriend: true,
		StatusUpdate: "",
	}
	sqlStatement:=`INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate") VALUES ($1,$2,$3,$4) RETURNING "RelationshipId"`
	err:=db.QueryRow(sqlStatement,userRelationship.UserEmail,userRelationship.FriendEmail,userRelationship.IsFriend,userRelationship.StatusUpdate).Scan(&id)
	if err!=nil{
		return err
	}
	return nil
}

func GetFriendListByEmail(db *sql.DB, email string) (*models.FriendsList, error) {
	response:=&models.FriendsList{}
	sqlStatement:=`select u."UserEmail" as Friends
								from "UserRelationship" u
								where (u."UserEmail"=$1 or u."FriendEmail"=$1)
								and "IsFriend"=true and u."UserEmail" not like $1
								Union
								select r."FriendEmail" as Friends
								from "UserRelationship" r
								where (r."UserEmail"=$1 or r."FriendEmail"=$1)
								and "IsFriend"=true and r."FriendEmail" not like $1`
	rows,err:=db.Query(sqlStatement,email)
	if err!=nil{
		return response, err
	}
	defer rows.Close()

	for rows.Next(){
		var emailTemp string
		err = rows.Scan(&emailTemp)
		if err!=nil{
			fmt.Println("email", err)
			return response, err
		}
		response.Friends = append(response.Friends, emailTemp)
	}

	return response, nil
}

func GetCommonFriend(db *sql.DB, userEmail string, friendEmail string) (*models.FriendsList,error) {
	responseFriend:=&models.FriendsList{}
	sqlStatement:=`select *
								from (select u."UserEmail"
								from public."UserRelationship" u
								where ("UserEmail"=$1 or "FriendEmail"=$1)
								and "IsFriend"=true and u."UserEmail" not like $1
								Union
								select u."FriendEmail" as Friends
								from public."UserRelationship" u
								where ("UserEmail"=$1 or "FriendEmail"=$1)
								and "IsFriend"=true and u."FriendEmail" not like $1) a
								where "UserEmail" in 
								(select r."UserEmail" as Friends
								from public."UserRelationship" r
								where (r."UserEmail"=$2 or r."FriendEmail"=$2)
								and r."IsFriend"=true and r."UserEmail" not like $2
								Union
								select s."FriendEmail" as Friends
								from public."UserRelationship" s
								where (s."UserEmail"=$2 or s."FriendEmail"=$2)
								and s."IsFriend"=true and s."FriendEmail" not like $2)`
	rows,err:=db.Query(sqlStatement,userEmail,friendEmail)
	if err!=nil{
		return responseFriend, err
	}
	defer rows.Close()
	for rows.Next(){
		var friend string
		err:=rows.Scan(&friend)
		if err!=nil{
			return responseFriend, err
		}
		responseFriend.Friends=append(responseFriend.Friends,friend)
	}
	return responseFriend, nil
}

func UpdateRelationship(db *sql.DB, requestor string, target string, status string) error {
	sqlUpdate:=`UPDATE "UserRelationship"
								SET "StatusUpdate"=$3 
								WHERE "UserEmail"=$1 and "FriendEmail"=$2`
	_,err:=db.Exec(sqlUpdate,requestor,target,status)
	if err!=nil{
		return err
	}
	return nil
}

func InsertRelationship(db *sql.DB, userEmail string, friendEmail string, status string) error {
	userRelationship:= models.UserRelationship{
		UserEmail: userEmail,
		FriendEmail: friendEmail,
		IsFriend: false,
		StatusUpdate: status,
	}
	var id int64
	sqlStatement:=`INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate") VALUES ($1,$2,$3,$4) RETURNING "RelationshipId"`
	err:=db.QueryRow(sqlStatement,userRelationship.UserEmail,userRelationship.FriendEmail,userRelationship.IsFriend,userRelationship.StatusUpdate).Scan(&id)
	if err!=nil{
		return err
	}
	return nil
}

func GetAllSubscriberByEmail(db *sql.DB, email string, text string) (*models.FriendsList, error) {
	subscribeResponse:=&models.FriendsList{}
	var emailSubscribeList []string
	var arrExtractEmail []string
	sqlStatement:=`select "UserEmail" as Friends
								from public."UserRelationship" u
								where (u."IsFriend"=true and u."StatusUpdate"!='BLOCK'
									   	and (u."UserEmail"=$1 
										or u."FriendEmail"=$1)
									  	and u."UserEmail" not like $1) or
										(u."StatusUpdate"='SUBSCRIBE' 
										and (u."UserEmail"=$1 
										or u."FriendEmail"=$1)
										and u."UserEmail" not like $1)
								union
								select u."FriendEmail" as Friends
								from public."UserRelationship" u
								where (u."IsFriend"=true and u."StatusUpdate"!='BLOCK'
									   	and (u."UserEmail"=$1 
										or u."FriendEmail"=$1)
									  	and u."FriendEmail" not like $1) or
										(u."StatusUpdate"='SUBSCRIBE' 
										and (u."UserEmail"=$1 
										or u."FriendEmail"=$1)
										and u."FriendEmail" not like $1)
										`
	rows,err:=db.Query(sqlStatement,email)
	if err!=nil{
		return subscribeResponse, err
	}
	defer rows.Close()
	for rows.Next(){
		var friend string
		err:=rows.Scan(&friend)
		if err!=nil{
			log.Println("Cannot scan to friend", err.Error())
			return subscribeResponse, err
		}
		emailSubscribeList=append(emailSubscribeList,friend)
	}
	if text != ""{
		var extractEmail = GetAllEmail(text)
		for i:=0;i<len(extractEmail);i++{
			if err:=CheckUserExistsInDb(db,extractEmail[i]); err!=nil{
				log.Println("Error ", err.Error())
				//return subscribeResponse, err
			}else {
				arrExtractEmail=append(arrExtractEmail,extractEmail[i])
			}
		}

	}
	emailsSubscribe := removeDuplicates(append([]string{}, append(emailSubscribeList, arrExtractEmail...)...))
	subscribeResponse.Friends=emailsSubscribe
	return subscribeResponse,nil
}

func GetAllEmail(text string) []string{
	re:=regexp.MustCompile(`[a-zA-Z0-9]+@[a-zA-Z0-9\.]+\.[a-zA-Z0-9]+`)
	match:=re.FindAllString(text,-1)
	return match
}

func removeDuplicates(elements []string) []string{
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range elements {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
