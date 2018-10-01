package repository

import (
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/modules/address/model"
)

type AddressRepositoryImpl struct {
	db *sql.DB
}

func NewAddressRepositoryImpl(db *sql.DB) *AddressRepositoryImpl {
	return &AddressRepositoryImpl{db}
}

func (r *AddressRepositoryImpl) Save(obj *model.Address) error{
	query := `INSERT INTO tb_alamat (id_alamat,id_user,alamat_lengkap,provinsi,kabupaten,kecamatan,kode_pos,created_time,updated_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.IdAlamat,obj.IdUser,obj.Alamat,obj.Provinsi,obj.Kabupaten,obj.Kecamatan,obj.KodePos,obj.CreatedAt,obj.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *AddressRepositoryImpl) Update(id string,obj *model.Address) error{
	query := `UPDATE tb_alamat SET alamat_lengkap= ?,provinsi= ?,kabupaten= ?,kecamatan= ?,kode_pos= ?, updated_time = ? where id_alamat = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.Alamat,obj.Provinsi,obj.Kabupaten,obj.Kecamatan,obj.KodePos,obj.UpdatedAt,id)
	if err != nil {
		return err
	}

	return nil
}

func (r *AddressRepositoryImpl) Delete(id string) error{
	query := `DELETE tb_alamat where id_alamat = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (r *AddressRepositoryImpl) FindByUserId(idUser string) (*model.Addresss,error){
	query := `SELECT id_alamat,IFNULL(alamat_lengkap,"") as alamat_lengkap,kabupaten,provinsi,
CONCAT(a.alamat_lengkap,' - ',k.nama,' - ',p.nama) as alamat FROM tb_alamat a
LEFT JOIN tb_provinsi p ON a.provinsi = p.id_provinsi
LEFT JOIN tb_kabupaten k ON a.kabupaten = k.id_kabupaten
WHERE id_user = "`+idUser+`"`

	var objs model.Addresss

	rows, err := r.db.Query(query)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.Address

		err = rows.Scan(&obj.IdAlamat,&obj.Alamat,&obj.Kabupaten,&obj.Provinsi,&obj.AlamatLengkap)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}