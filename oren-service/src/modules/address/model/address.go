package model

import "time"

type Address struct {
	IdAlamat string
	IdUser string
	Alamat string
	AlamatLengkap string
	Provinsi string
	Kabupaten string
	Kecamatan string
	KodePos string
	CreatedAt int64
	UpdatedAt int64
}

// List User
type Addresss []Address

//  Constructor
func NewAddress() *Address{
	return &Address{
		CreatedAt:time.Now().Unix(),
		UpdatedAt:time.Now().Unix(),
	}
}
