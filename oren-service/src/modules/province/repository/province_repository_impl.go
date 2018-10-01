package repository

import (
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/modules/province/model"
)

type ProvinceRepositoryImpl struct {
	db *sql.DB
}

func NewProvinceRepositoryImpl(db *sql.DB) *ProvinceRepositoryImpl {
	return &ProvinceRepositoryImpl{db}
}

func (r *ProvinceRepositoryImpl) FindAll() (*model.Provinces,error){
	query := `SELECT * FROM tb_provinsi`

	var objs model.Provinces

	rows, err := r.db.Query(query)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.Province

		err = rows.Scan(&obj.IdProvinsi,&obj.Nama)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}

func (r *ProvinceRepositoryImpl) FindAllCity() (*model.Citys,error){
	query := `SELECT * FROM tb_kabupaten`

	var objs model.Citys

	rows, err := r.db.Query(query)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.City

		err = rows.Scan(&obj.IdKabupaten,&obj.IdProvinsi,&obj.Nama)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}