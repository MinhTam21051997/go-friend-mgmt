package user

import (
	"database/sql"
	"errors"
	"go-friend-mgmt/cmd/internal/services/models"
	"go-friend-mgmt/cmd/internal/services/utils"
)

func (s ServiceImpl)CreateUser(email models.EmailUser) (*models.Response, error)  {
	response:=&models.Response{}
	err:=utils.VaditationEmail(email.Email)
	if err!=nil{
		return response,err
	}
	count, err:=CheckUserExists(s.DB, email.Email)
	if err!=nil{
		return  response, err
	}
	if count>=1{
		return response, errors.New("Email already exists in db")
	}
	err=CreateUserByEmail(s.DB, email.Email)
	if err!=nil{
		return response, err
	}
	response.Success=true
	return response,nil
}

func CreateUserByEmail(db *sql.DB, emailUser string) error {
	sqlStatement:=`INSERT INTO "User" ("email") VALUES ($1)`
	_,err:=db.Exec(sqlStatement,emailUser)
	if err!=nil{
		return err
	}
	return nil
}

func CheckUserExists(db *sql.DB, email string) (int,error) {
	var count int
	err:=utils.VaditationEmail(email)
	if err!=nil{
		return count,err
	}else {
		err = db.QueryRow(`SELECT COUNT(*) FROM "User" WHERE "email" = $1`, email).Scan(&count)
		if err!=nil{
			return count,err
		}
	}
	return count,nil
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
