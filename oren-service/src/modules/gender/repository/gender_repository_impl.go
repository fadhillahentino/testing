package repository

import (
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/modules/gender/model"
)

type GenderRepositoryImpl struct {
	db *sql.DB
}

func NewGenderRepositoryImpl(db *sql.DB) *GenderRepositoryImpl {
	return &GenderRepositoryImpl{db}
}

func (r *GenderRepositoryImpl) FindAll() (*model.Genders,error){
	query := `SELECT * FROM tb_gender`

	var objs model.Genders

	rows, err := r.db.Query(query)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.Gender

		err = rows.Scan(&obj.IdGender,&obj.NamaGender)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}










