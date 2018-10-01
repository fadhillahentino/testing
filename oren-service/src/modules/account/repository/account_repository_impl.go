package repository

import (
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/modules/account/model"
)

type AccountRepositoryImpl struct {
	db *sql.DB
}

func NewAccountRepositoryImpl(db *sql.DB) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{db}
}

func (r *AccountRepositoryImpl) Save(obj *model.Account) error{
	query := `INSERT INTO tb_rekening (id_rekening,nomor_rekening,nama_bank,atas_nama,id_user) VALUES (?, ?, ?, ?, ?)`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.IdRekening,obj.NomorRekening,obj.NamaBank,obj.AtasNama,obj.IdUser)
	if err != nil {
		return err
	}

	return nil
}

func (r *AccountRepositoryImpl) Update(id string,obj *model.Account) error{
	query := `UPDATE tb_rekening SET nomor_rekening = ?, nama_bank = ?, atas_nama = ? where id_rekening = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.NomorRekening,obj.NamaBank,obj.AtasNama, obj.IdRekening)
	if err != nil {
		return err
	}

	return nil
}

func (r *AccountRepositoryImpl) FindById(IdUser string) (*model.Accounts,error){
	query := `SELECT id_rekening,nomor_rekening,nama_bank,atas_nama,id_user,CONCAT(nomor_rekening,' - ',nama_bank,' - ',atas_nama) as rekening
FROM tb_rekening WHERE id_user = '`+IdUser+`';`

	var objs model.Accounts

	rows, err := r.db.Query(query)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.Account

		err = rows.Scan(&obj.IdRekening,&obj.NomorRekening,&obj.NamaBank,&obj.AtasNama,&obj.IdUser,&obj.Rekening)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}








