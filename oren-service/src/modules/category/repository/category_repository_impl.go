package repository

import (
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/modules/category/model"
	"fmt"
)

type CategoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepositoryImpl(db *sql.DB) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{db}
}

func (r *CategoryRepositoryImpl) FindbyGender(idGender string) (*model.Categorys,error){
	query := fmt.Sprintf(`SELECT id_kategori,nama_kategori FROM tb_kategori where id_gender = '%s'`,idGender)

	var objs model.Categorys

	rows, err := r.db.Query(query)
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	for rows.Next(){
		var obj model.Category

		err = rows.Scan(&obj.IdCategory,&obj.NamaCategory)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}

func (r *CategoryRepositoryImpl) FindAll() (*model.Categorys,error){
	query := `SELECT * FROM tb_kategori`

	var objs model.Categorys

	rows, err := r.db.Query(query)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.Category

		err = rows.Scan(&obj.IdCategory,&obj.NamaCategory,&obj.IdGender)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}










