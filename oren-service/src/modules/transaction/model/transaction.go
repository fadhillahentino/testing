package model

import "time"

type Transaction struct {
	IdBusana string
	IdTransaksi string
	IdPenyewa string
	IdPemilik string
	StartDate int64
	EndDate int64
	TotalHarga float64
	BuktiTransfer string
	BuktiPengiriman string
	BuktiPengirimanKembali string
	UrlBuktiTransfer string
	UrlBuktiPengiriman string
	UrlBuktiPengirimanKembali string
	IdPengiriman string
	IdPengirimanKembali string
	IdStatus string
	IdAlamat string
	CreatedAt int64
	UpdatedAt int64
	NomorRekening string
	NamaBank string
	AtasNama string
	NominaTransfer float64
	Resi string
	IdAlamatPemilik string
	HargaPengiriman float64
	IdRekening string
}

type Transactions []Transaction

//  Constructor
func NewTransaction() *Transaction{
	return &Transaction{
		CreatedAt:time.Now().Unix(),
		UpdatedAt:time.Now().Unix(),
	}
}

type Notification struct {
	IdTransaksi string
	NamaBusana string
	Penyewa string
	Pemilik string
	StartDate string
	EndDate string
	IdStatus string
	Status string
	UpdatedAt int64
	Flag string
}

type Notifications []Notification

//  Constructor
func NewNotification() *Notification{
	return &Notification{}
}

type RentDetail struct {
	IdUserPenyewa string
	IdUserPemilik string
	IdTransaksi string
	NamaBusana string
	IdStatus string
	Penyewa string
	Pemilik string
	StartDate string
	EndDate string
	RatingPenyewa float64
	RatingPemilik float64
	UrlPenyewa string
	UrlPemilik string
	UrlBusana string
	Pengiriman string
	Harga float64
	HargaPengiriman float64
	Deposit float64
	HargaSewa float64
	TotalHarga float64
	Alamat string
	Kecamatan string
	Kabupaten string
	Provinsi string
	KodePos string
	NoHp string
	UrlBuktiTransfer string
	NomorRekening string
	NamaBank string
	AtasNama string
	NominaTransfer float64
	BuktiPengiriman string
	BuktiPengirimanKembali string
	UrlBuktiPengiriman string
	UrlBuktiPengirimanKembali string
	IdUser string
}

type RentDetails []RentDetail

//  Constructor
func NewRentDetail() *RentDetail{
	return &RentDetail{}
}