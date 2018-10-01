package model

type Category struct {
	IdCategory string
	NamaCategory string
	IdGender string
}

type Categorys []Category

func NewModel() *Category{
	return &Category{}
}