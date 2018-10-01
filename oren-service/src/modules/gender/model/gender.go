package model

type Gender struct {
	IdGender string
	NamaGender string
}

type Genders []Gender

func NewModel() *Gender{
	return &Gender{}
}