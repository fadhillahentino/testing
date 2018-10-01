package model

import "time"

// User Struct
type User struct {
	IdUser string
	Phone string
	Name string
	Email string
	Password string
	Foto string
	UrlFoto string
	Uid string
	CreatedAt int64
	UpdatedAt int64
}

// List User
type Users []User

//  Constructor
func NewUser() *User{
	return &User{
		CreatedAt:time.Now().Unix(),
		UpdatedAt:time.Now().Unix(),
	}
}