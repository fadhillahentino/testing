package model

type Province struct {
	IdProvinsi string
	Nama string
}

type Provinces []Province

//  Constructor
func NewProvince() *Province{
	return &Province{
	}
}

type City struct {
	IdKabupaten string
	IdProvinsi string
	Nama string
}

type Citys []City

//  Constructor
func NewCity() *City{
	return &City{
	}
}