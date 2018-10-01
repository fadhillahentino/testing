package model

type Shippment struct {
	IdPengiriman string
	Pengiriman string
}

type Shippments []Shippment

//  Constructor
func NewShippment() *Shippment{
	return &Shippment{}
}

type ShippmentFashion struct {
	IdPengiriman string
	IdBusana string
}

type ShippmentFashions []ShippmentFashion