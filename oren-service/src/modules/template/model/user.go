package model

import "time"

type User struct {
	IdUser string
	NoHp string
	Nama string
	Email string
	Password string
	Foto string
	CreatedAt int64
	UpdatedAt int64
}

type Users []User

//  Constructor
func NewUser() *User{
	return &User{
		CreatedAt:time.Now().Unix(),
		UpdatedAt:time.Now().Unix(),
	}
}