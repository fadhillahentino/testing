package model

import "time"

type Account struct {
	IdRekening string
	NomorRekening string
	NamaBank string
	AtasNama string
	Rekening string
	IdUser string
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

type Accounts []Account

//  Constructor
func NewAccount() *Account{
	return &Account{
		CreatedAt:time.Now().Unix(),
		UpdatedAt:time.Now().Unix(),
		DeletedAt:time.Now().Unix(),
	}
}