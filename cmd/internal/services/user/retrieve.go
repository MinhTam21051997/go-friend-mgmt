package user

import (
	"database/sql"
	"go-friend-mgmt/cmd/internal/services/models"
	"go-friend-mgmt/cmd/internal/services/utils"
	"log"
	"regexp"
)

func (ServiceImpl) RetrieveByID(ID int) (*models.Users, error) {
	return &models.Users{
		ID:    1,
		Name:  "john doe",
		Email: "john.doe@email.com",
	}, nil
}

func (s ServiceImpl) CreateConnectionFriend(friend models.Friend) (int64 ,*models.Response ) {
	//db:= database.ConnectionDB()
	//defer db.Close()

	var id int64

	if friend.UserEmail=="" && friend.FriendEmail==""{
		return 0, &models.Response{
			Success: false,
			Message: "Please input email!",
		}
	}else{
		isValidUserEmail:=utils.VaditationEmail(friend.UserEmail)
		isValidFriendEmail:=utils.VaditationEmail(friend.FriendEmail)
		if isValidUserEmail==false || isValidFriendEmail==false{
			return 0,&models.Response{
				Success: false,
				Message: "Invalid email!",
			}
		}else {
			isCheckUserEmailInDb:=CheckUserExists(s.DB,friend.UserEmail)
			isCheckFriendEmailInDb:=CheckUserExists(s.DB,friend.FriendEmail)
			if isCheckUserEmailInDb==false || isCheckFriendEmailInDb==false{
				return 0,&models.Response{
					Success: false,
					Message: "Email does not exist in database!",
				}
			}else {
				isFriend:=CheckIsFriendOrBlock(s.DB,friend.UserEmail,friend.FriendEmail)
				log.Println("isFriend in create connection friend", isFriend)
				if isFriend{
					return 0, &models.Response{
						Success: false,
						Message: "They were blocked or were friends !",
					}
				}else {
					userRelationship:= models.UserRelationship{
						UserEmail: friend.UserEmail,
						FriendEmail: friend.FriendEmail,
						IsFriend: true,
						StatusUpdate: "",
					}
					log.Println("userRelationship  " , userRelationship)
					sqlStatement:=`INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate") VALUES ($1,$2,$3,$4) RETURNING "RelationshipId"`
					err:=s.DB.QueryRow(sqlStatement,userRelationship.UserEmail,userRelationship.FriendEmail,userRelationship.IsFriend,userRelationship.StatusUpdate).Scan(&id)

					if err!=nil{
						log.Println("err " , err.Error())
						return 0, &models.Response{
							Success: false,
							Message: "Create connection friends failed !",
						}
					}else {
						return id,&models.Response{
							Success: true,
							Message: "Created connection friends successfully !",
						}
					}
				}
			}
		}
	}

	return id,&models.Response{
		Success: false,
		Message: "Create connection friends failed !",
	}
}

func (s ServiceImpl) ReceiveFriendListByEmail(email string) (*models.ResponseFriend, error)  {
	//db:= database.ConnectionDB()
	//defer db.Close()

	responseFriend:= &models.ResponseFriend{
		Success: false,
		Friends: nil,
		Count: 0,
		Message: "",
	}

	if email==""{
		responseFriend.Message="Email cannot be empty!"
		return responseFriend,nil
	}else {
		isValidEmail:=utils.VaditationEmail(email)
		if !isValidEmail{
			responseFriend.Message="Invalid email!"
			return responseFriend,nil
		}else {
			isCheckEmailInDb:=CheckUserExists(s.DB,email)
			if !isCheckEmailInDb{
				responseFriend.Message="Email does not exist in database!"
				return responseFriend,nil
			}else{
				sqlStatement:=`select u."UserEmail" as Friends
								from "UserRelationship" u
								where (u."UserEmail"=$1 or u."FriendEmail"=$1)
								and "IsFriend"=true and u."UserEmail" not like $1
								Union
								select r."FriendEmail" as Friends
								from "UserRelationship" r
								where (r."UserEmail"=$1 or r."FriendEmail"=$1)
								and "IsFriend"=true and r."FriendEmail" not like $1`
				rows,err:=s.DB.Query(sqlStatement,email)
				if err!=nil{
					log.Println("[ReceiveFriendList]: ", err.Error())
					responseFriend.Message="Cannot query on table user_relationship"
					return responseFriend, err
				}
				defer rows.Close()

				for rows.Next(){
					var emailTemp string
					err = rows.Scan(&emailTemp)
					if err!=nil{
						log.Fatalf("Cannot insert email to slice %v ", err.Error() )
					}
					responseFriend.Friends = append(responseFriend.Friends, emailTemp)
					responseFriend.Success=true
					responseFriend.Count=responseFriend.Count+1
					responseFriend.Message="Get receive friends list successfully!"
				}
			}
		}
	}

	return responseFriend, nil
}

func (s ServiceImpl) ReceiveCommonFriendList(friend models.Friend) (*models.ResponseFriend,error)  {
	//db:= database.ConnectionDB()

	//defer db.Close()

	responseFriend:= &models.ResponseFriend{
		Success: false,
		Friends: nil,
		Count: 0,
		Message: "",
	}

	if friend.UserEmail=="" && friend.FriendEmail==""{
		responseFriend.Message="Email cannot be empty!"
		return responseFriend,nil
	}else{
		isValidUserEmail:=utils.VaditationEmail(friend.UserEmail)
		isValidFriendEmail:=utils.VaditationEmail(friend.FriendEmail)
		if isValidUserEmail==false || isValidFriendEmail==false{
			responseFriend.Message="Invalid email!"
			return responseFriend,nil
		}else {
			isCheckUserEmailInDb:=CheckUserExists(s.DB,friend.UserEmail)
			isCheckFriendEmailInDb:=CheckUserExists(s.DB,friend.FriendEmail)
			if isCheckUserEmailInDb==false || isCheckFriendEmailInDb==false{
				responseFriend.Message="Email does not exist in database!"
				return responseFriend,nil
			}else {
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
				rows,err:=s.DB.Query(sqlStatement,friend.UserEmail,friend.FriendEmail)
				if err!=nil{
					log.Println("Cannot query get common friends list :", err.Error())
					responseFriend.Message="Cannot query get common friends list"
					return responseFriend, err
				}else {
					defer rows.Close()
					log.Println("rows :", rows)
					var commonFriendList []string
					for rows.Next(){
						var friend string
						err:=rows.Scan(&friend)
						if err!=nil{
							log.Println("Cannot scan to friend", err.Error())
						}else {
							commonFriendList=append(commonFriendList,friend)
						}
					}
					log.Println("commonFriendList :", commonFriendList)
					responseFriend.Success=true
					responseFriend.Friends=commonFriendList
					responseFriend.Count=int64(len(commonFriendList))
					responseFriend.Message="Get common friends list between two email successfully !"
				}
			}
		}
	}

	return responseFriend, nil
}

func (s ServiceImpl) SubscribeUpdateFromEmail(subscribeUser models.SubscribeUser) (*models.Response,error)  {
	//db:= database.ConnectionDB()

	//defer db.Close()

	response:= &models.Response{
		Success: false,
		Message: "",
	}

	if subscribeUser.Requestor=="" && subscribeUser.Target==""{
		response.Message="Email cannot be empty!"
		return response,nil
	}else{
		if subscribeUser.Requestor==subscribeUser.Target{
			response.Message="The two emails cannot be the same!"
			return response,nil
		}else {
			isValidUserEmail:=utils.VaditationEmail(subscribeUser.Requestor)
			isValidFriendEmail:=utils.VaditationEmail(subscribeUser.Target)
			if isValidUserEmail==false || isValidFriendEmail==false{
				response.Message="Invalid email!"
				return response,nil
			}else {
				isCheckUserEmailInDb:=CheckUserExists(s.DB,subscribeUser.Requestor)
				isCheckFriendEmailInDb:=CheckUserExists(s.DB,subscribeUser.Target)
				if isCheckUserEmailInDb==false || isCheckFriendEmailInDb==false{
					response.Message="Email does not exist in database!"
					return response,nil
				}else {
					isRelationship:=CheckRelationshipExists(s.DB,subscribeUser.Requestor, subscribeUser.Target)
					if isRelationship{
						sqlUpdate:=`UPDATE "UserRelationship"
								SET "StatusUpdate"=$3 
								WHERE "UserEmail"=$1 and "FriendEmail"=$2`
						_,err:=s.DB.Exec(sqlUpdate,subscribeUser.Requestor,subscribeUser.Target,"SUBSCRIBE")
						if err!=nil{
							log.Println("Cannot subcribe to update ", err.Error())
							response.Message="Cannot query subscribe to update"
							return response, err
						}else {
							response.Success=true
							response.Message="Subscribe to updates successfully!"
							return response, nil
						}
					}else {
						userRelationship:= models.UserRelationship{
							UserEmail: subscribeUser.Requestor,
							FriendEmail: subscribeUser.Target,
							IsFriend: false,
							StatusUpdate: "SUBSCRIBE",
						}

						var id int64
						sqlStatement:=`INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate") VALUES ($1,$2,$3,$4) RETURNING "RelationshipId"`

						err:=s.DB.QueryRow(sqlStatement,userRelationship.UserEmail,userRelationship.FriendEmail,userRelationship.IsFriend,userRelationship.StatusUpdate).Scan(&id)

						if err!=nil{
							log.Println("Cannot insert subscribe to updates " , err.Error())
							response.Message="Cannot insert subscribe to updates"
							return response,err
						}else {
							response.Success=true
							response.Message="Subscribe to updates successfully!"
							return response,nil
						}
					}
				}
			}
		}
	}

	return response, nil
}

func (s ServiceImpl) BlockUpdateFromEmail(subscribeUser models.SubscribeUser) (*models.Response,error)  {
	//db:= database.ConnectionDB()

	//defer db.Close()

	response:= &models.Response{
		Success: false,
		Message: "",
	}

	if subscribeUser.Requestor=="" && subscribeUser.Target==""{
		response.Message="Email cannot be empty!"
		return response,nil
	}else{
		if subscribeUser.Requestor==subscribeUser.Target{
			response.Message="The two emails cannot be the same!"
			return response,nil
		}else {
			isValidUserEmail:=utils.VaditationEmail(subscribeUser.Requestor)
			isValidFriendEmail:=utils.VaditationEmail(subscribeUser.Target)
			if isValidUserEmail==false || isValidFriendEmail==false{
				response.Message="Invalid email!"
				return response,nil
			}else {
				isCheckUserEmailInDb:=CheckUserExists(s.DB,subscribeUser.Requestor)
				isCheckFriendEmailInDb:=CheckUserExists(s.DB,subscribeUser.Target)
				if isCheckUserEmailInDb==false || isCheckFriendEmailInDb==false{
					response.Message="Email does not exist in database!"
					return response,nil
				}else {
					isRelationship:=CheckRelationshipExists(s.DB,subscribeUser.Requestor, subscribeUser.Target)
					if isRelationship{
						sqlUpdate:=`UPDATE "UserRelationship"
								SET "StatusUpdate"=$3
								WHERE "UserEmail"=$1 and "FriendEmail"=$2`
						_,err:=s.DB.Exec(sqlUpdate,subscribeUser.Requestor,subscribeUser.Target,"BLOCK")
						if err!=nil{
							log.Println("Cannot update block to db " , err.Error())
							response.Success=false
							response.Message="Block to updates failed!"
							return response, nil
						}else {
							response.Success=true
							response.Message="Block to updates successfully!"
							return response, nil
						}

					}else {
						userRelationship:= models.UserRelationship{
							UserEmail: subscribeUser.Requestor,
							FriendEmail: subscribeUser.Target,
							IsFriend: false,
							StatusUpdate: "BLOCK",
						}

						var id int64
						sqlStatement:=`INSERT INTO "UserRelationship" ("UserEmail","FriendEmail","IsFriend","StatusUpdate") VALUES ($1,$2,$3,$4) RETURNING "RelationshipId"`
						err:=s.DB.QueryRow(sqlStatement,userRelationship.UserEmail,userRelationship.FriendEmail,userRelationship.IsFriend,userRelationship.StatusUpdate).Scan(&id)

						if err!=nil{
							log.Println("Cannot insert block to updates " , err.Error())
							response.Message="Cannot insert block to updates"
							return response,err
						}else {
							response.Success=true
							response.Message="Block to updates successfully!"
							return response,nil
						}
					}
				}
			}
		}
	}

	return response, nil
}

func (s ServiceImpl) GetAllSubscribeUpdateByEmail(retrieve models.RetrieveUpdate) (*models.SubscribeResponse,error)  {
	subscribeResponse:= &models.SubscribeResponse{
		Success: false,
		Recipients: nil,
		Message: "",
	}
	var emailSubscribeList []string
	var arrExtractEmail []string

	if retrieve.Sender==""{
		subscribeResponse.Message="Email cannot be empty!"
		return subscribeResponse,nil
	}else{
		isValidSenderEmail:=utils.VaditationEmail(retrieve.Sender)
		if isValidSenderEmail==false{
			subscribeResponse.Message="Invalid email!"
			return subscribeResponse,nil
		}else {
			isCheckSenderEmailInDb:=CheckUserExists(s.DB,retrieve.Sender)
			if isCheckSenderEmailInDb==false{
				subscribeResponse.Message="Email does not exist in database!"
				return subscribeResponse,nil
			}else {
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
				rows,err:=s.DB.Query(sqlStatement,retrieve.Sender)
				if err!=nil{
					log.Println("Cannot query get email subscribe :", err.Error())
					if retrieve.Text != ""{
						var subEmail = GetAllEmail(retrieve.Text)
						log.Println("extractEmail :", subEmail)
						for i:=0;i<len(subEmail);i++{
							if CheckUserExists(s.DB,subEmail[i]){
								arrExtractEmail=append(arrExtractEmail,subEmail[i])
							}
						}
					}
					emailsSubscribe := removeDuplicates(append([]string{}, append(emailSubscribeList, arrExtractEmail...)...))
					log.Println("emailsSubscribe :", emailsSubscribe)
					subscribeResponse.Success=true
					subscribeResponse.Recipients=emailsSubscribe
					//subscribeResponse.Message="Cannot query get email subscribe"
					subscribeResponse.Message="Get list email subscribe by email successfully !"
					return subscribeResponse, err
				}else {
					defer rows.Close()
					for rows.Next(){
						var friend string
						err:=rows.Scan(&friend)
						if err!=nil{
							log.Println("Cannot scan to friend", err.Error())
						}else {
							emailSubscribeList=append(emailSubscribeList,friend)
						}
					}
					log.Println("Email Subscribe List :", emailSubscribeList)
					if retrieve.Text != ""{
						var extractEmail = GetAllEmail(retrieve.Text)
						log.Println("extractEmail :", extractEmail)
						for i:=0;i<len(extractEmail);i++{
							if CheckUserExists(s.DB,extractEmail[i]){
								arrExtractEmail=append(arrExtractEmail,extractEmail[i])
							}
						}

					}
					emailsSubscribe := removeDuplicates(append([]string{}, append(emailSubscribeList, arrExtractEmail...)...))
					log.Println("emailsSubscribe :", emailsSubscribe)
					subscribeResponse.Success=true
					subscribeResponse.Recipients=emailsSubscribe
					subscribeResponse.Message="Get list email subscribe by email successfully !"
				}
			}
		}
	}

	return subscribeResponse, nil
}

func CheckUserExists(db *sql.DB, email string) bool {
	//db:= database.ConnectionDB()
	//defer db.Close()

	var err error
	var count int64
	if email==""{
		log.Println("Email cannot be empty!")
		return false
	}else {
		err = db.QueryRow(`SELECT COUNT(*) FROM "User" WHERE "email" = $1`, email).Scan(&count)
		if err!=nil{
			log.Println("[CheckExitsEmail] query error : ", err.Error())
			return false
		}else {
			if count==1{
				return true
			}
		}
	}
	return false
}

func CheckRelationshipExists(db *sql.DB, userEmail string, friendEmail string) bool {
	//db:= database.ConnectionDB()
	//defer db.Close()

	var err error
	var count int64
	if userEmail=="" || friendEmail==""{
		log.Println("Email cannot be empty!")
		return false
	}else {
		err = db.QueryRow(`SELECT COUNT(*) FROM "UserRelationship" WHERE "UserEmail" = $1 and "FriendEmail" = $2`, userEmail, friendEmail).Scan(&count)
		if err!=nil{
			log.Println("[CheckRelationshipExists] query error : ", err.Error())
			return false
		}else {
			if count>=1{
				return true
			}
		}
	}
	return false
}

func CheckIsFriendOrBlock(db *sql.DB, userEmail string, friendEmail string) bool {
	isFriend:=false
	if userEmail=="" || friendEmail==""{
		log.Println("[CheckIsFriendOrBlock] Email cannot be empty!")
		return false
	}else {
		sqlQuery:=`SELECT * FROM "UserRelationship" WHERE ("UserEmail" = $1 and "FriendEmail" = $2) or ("FriendEmail" = $1 and "UserEmail" = $2)`
		rows,err := db.Query(sqlQuery, userEmail, friendEmail)
		if err!=nil{
			log.Println("[CheckIsFriendOrBlock] query error : ", err.Error())
			isFriend=true
			return isFriend
		}
		defer rows.Close()
		var relationships []models.UserRelationship
		for rows.Next(){
			var rel models.UserRelationship
			err = rows.Scan(&rel.RelationshipId,&rel.UserEmail,&rel.FriendEmail,&rel.IsFriend,&rel.StatusUpdate)
			if err!=nil{
				log.Fatalf("[CheckIsFriendOrBlock] Cannot scan to UserRelationship model %v ", err.Error() )
			}
			relationships=append(relationships,rel)
		}

		log.Println("len(relationships) in create connection friend", len(relationships))
		if len(relationships)<1{
			isFriend=false
			return isFriend
		}else {
			for i:=0;i<len(relationships);i++{
				if relationships[i].StatusUpdate=="BLOCK" || relationships[i].IsFriend==true{
					isFriend=true
					return isFriend
				}
			}
		}
	}
	return isFriend
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
