package model

type Rating struct {
	IdRating string
	RespChat float64
	Flexible float64
	Perawatan float64
	TepatWaktu float64
	Friendly float64
	Total float64
	IdUser string
	Desc string
	IdTransaksi string
	Flag int64
	Nama string
	JmlData int
	Status string
}

type Ratings []Rating

//  Constructor
func NewRating() *Rating{
	return &Rating{}
}

type RatingSummary struct {
	RespChat float64
	Flexible float64
	Perawatan float64
	TepatWaktu float64
	Friendly float64
	JmlData int64
	Total float64
	RespChatRat float64
	FlexibleRat float64
	PerawatanRat float64
	TepatWaktuRat float64
	FriendlyRat float64
	TotalRat float64
	IdUser string
}

type RatingSummarys []RatingSummary

//  Constructor
func NewRatingSummary() *RatingSummary{
	return &RatingSummary{}
}