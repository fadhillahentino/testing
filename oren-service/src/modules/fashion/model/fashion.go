package model

import "time"

// Struct
type Fashion struct {
	IdBusana string
	NamaBusana string
	IdKategori string
	IdUser string
	IdStatus string
	Berat float64
	Harga float64
	Deposit float64
	Deskripsi string
	FotoUtama string
	FotoSatu string
	FotoDua string
	FotoTiga string
	FotoEmpat string
	UrlFotoUtama string
	UrlFotoSatu string
	UrlFotoDua string
	UrlFotoTiga string
	UrlFotoEmpat string
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

// List
type Fashions []Fashion

//  Constructor
func NewFashion() *Fashion{
	return &Fashion{
		CreatedAt:time.Now().Unix(),
		UpdatedAt:time.Now().Unix(),
		DeletedAt:time.Now().Unix(),
	}
}

type FashionSearch struct {
	Id string
	IdStatus string
	Location string
	Rating float64
	Name string
	UrlImage string
	Price float64
}

type FashionSearchs []FashionSearch

type FashionDetail struct {
	IdBusana string
	NamaBusana string
	Berat float64
	Harga float64
	Deposit float64
	Deskripsi string
	UrlFotoUtama string
	UrlFotoSatu string
	UrlFotoDua string
	UrlFotoTiga string
	UrlFotoEmpat string
	IdUser string
	Nama string
	Phone string
	Alamat string
	Provinsi string
	Kabupaten string
	NamaKategori string
	Rating float64
	Pengiriman string
	Email string
	Password string
	Uid string
	UrlFotoUser string
}

func NewFashionDetail() *FashionDetail{
	return &FashionDetail{}
}