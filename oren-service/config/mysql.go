package config

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)

const user = "root"
//const user = "sa"
const pass = "123456"
//const url = "128.199.129.191"
const url = "localhost"
const port = "3306"
const db = "db_oren"

func GetMysqlDB()(*sql.DB,error){
	Oke()
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true",user,pass,url,port,db)
	db, err := sql.Open("mysql",connStr)
	if err != nil{
		fmt.Println(err.Error())
		return nil,err
	}
	db.SetMaxIdleConns(10000)
	db.SetMaxOpenConns(10000)
	Oke()
	return db, nil
}

func Oke()string{
	return "oke"
}