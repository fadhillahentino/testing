package repository

import (
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/modules/fashion/model"
	constant "github.com/fadhilfcr/oren-service/src/util"
	"time"
	"fmt"
)

const TBL_NAME = constant.TAG_BUSANA_TBL
const PK_ID = constant.TAG_BUSANA_PK

type FashionRepositoryImpl struct {
	db *sql.DB
}

func NewFashionRepositoryImpl(db *sql.DB) *FashionRepositoryImpl {
	return &FashionRepositoryImpl{db}
}

func (r *FashionRepositoryImpl) Save(obj *model.Fashion) error{
	query :=`INSERT INTO %s (id_busana, nama_busana, id_kategori, id_user, berat, harga, deposit, deskripsi, foto_utama, foto_satu, foto_dua, foto_tiga, foto_empat, id_status, created_time, updated_time, deleted_time, url_foto_utama, url_foto_satu, url_foto_dua, url_foto_tiga, url_foto_empat)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	querys :=  fmt.Sprintf(query,TBL_NAME)
	statement, err := r.db.Prepare(querys)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.IdBusana,obj.NamaBusana,obj.IdKategori,obj.IdUser,obj.Berat,obj.Harga,obj.Deposit,obj.Deskripsi,obj.FotoUtama,obj.FotoSatu,obj.FotoDua,obj.FotoTiga,obj.FotoEmpat,obj.IdStatus,obj.CreatedAt,obj.UpdatedAt,obj.DeletedAt,obj.UrlFotoUtama,obj.UrlFotoSatu,obj.UrlFotoDua,obj.UrlFotoTiga,obj.UrlFotoEmpat)
	if err != nil {
		return err
	}

	return nil
}

func (r *FashionRepositoryImpl) Update(id string,obj *model.Fashion)error{
	query := `UPDATE %s SET nama_busana = ?, id_kategori = ?, id_user = ?, berat = ?, harga = ?, deposit = ?, deskripsi = ?, foto_utama = ?, foto_satu = ?,
foto_dua = ?, foto_tiga = ?, foto_empat = ?, id_status = ?, updated_time = ?) where %s = ?`
	querys := fmt.Sprintf(query,TBL_NAME,PK_ID)
	statement, err := r.db.Prepare(querys)
	if err != nil {
		return err
	}

	defer statement.Close()
	obj.UpdatedAt = time.Now().Unix()
	_, err = statement.Exec(obj.NamaBusana,obj.IdKategori,obj.IdUser,obj.Berat,obj.Harga,obj.Deskripsi,obj.FotoUtama,obj.FotoSatu,obj.FotoDua,obj.FotoTiga,obj.FotoEmpat,obj.IdStatus,obj.UpdatedAt,obj.IdBusana)
	if err != nil {
		return err
	}

	return nil
}

func (r *FashionRepositoryImpl) UpdateStatus(id string,idStatus string)error{
	query := `UPDATE %s SET id_status = ?, updated_time = ? where %s = ?`
	querys := fmt.Sprintf(query,TBL_NAME,PK_ID)
	statement, err := r.db.Prepare(querys)
	if err != nil {
		return err
	}

	defer statement.Close()
	updTime := time.Now().Unix()
	_, err = statement.Exec(idStatus,updTime,id)
	if err != nil {
		return err
	}

	return nil
}

func (r *FashionRepositoryImpl) Delete(id string) error{
	query := `DELETE %s where %s = ?`
	querys := fmt.Sprintf(query,TBL_NAME,PK_ID)
	statement, err := r.db.Prepare(querys)
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

func (r *FashionRepositoryImpl) FindById(id string) (*model.FashionDetail,error){
	querys := `SELECT b.id_busana,b.nama_busana,b.berat,b.harga,b.deposit,b.deskripsi,b.url_foto_utama,b.url_foto_satu,b.url_foto_dua,b.url_foto_tiga,b.url_foto_empat,u.id_user,u.nama as name,u.no_hp,a.alamat_lengkap,ka.nama,p.nama,k.nama_kategori,IFNULL(rs.rating_total,0) as rating,u.email,u.password,u.uid,IFNULL(u.url_foto_user,'') as url_foto_user FROM tb_busana b
LEFT JOIN tb_user u ON b.id_user = u.id_user
LEFT JOIN tb_alamat a ON u.id_user = a.id_user
LEFT JOIN tb_rating_summary rs ON u.id_user = rs.id_user
LEFT JOIN tb_kategori k ON b.id_kategori = k.id_kategori
LEFT JOIN tb_kabupaten ka ON ka.id_kabupaten = a.kabupaten
LEFT JOIN tb_provinsi p ON p.id_provinsi = a.provinsi
WHERE b.id_busana = ?`
	var obj model.FashionDetail

	statement, err := r.db.Prepare(querys)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow(id).Scan(&obj.IdBusana,&obj.NamaBusana,&obj.Berat,&obj.Harga,&obj.Deposit,&obj.Deskripsi,&obj.UrlFotoUtama,&obj.UrlFotoSatu,&obj.UrlFotoDua,&obj.UrlFotoTiga,&obj.UrlFotoEmpat,&obj.IdUser,&obj.Nama,&obj.Phone,&obj.Alamat,&obj.Provinsi,&obj.Kabupaten,&obj.NamaKategori,&obj.Rating,&obj.Email,&obj.Password,&obj.Uid,&obj.UrlFotoUser)
	if err != nil {
		return nil,err
	}

	return &obj,nil
}

func (r *FashionRepositoryImpl) FindAll() (*model.Fashions,error){
	query := `SELECT * FROM %s`
	querys := fmt.Sprintf(query,TBL_NAME)
	var objs model.Fashions

	rows, err := r.db.Query(querys)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.Fashion

		err = rows.Scan(&obj.IdBusana,&obj.NamaBusana,&obj.IdKategori,&obj.IdUser,&obj.Berat,&obj.Harga,&obj.Deskripsi,&obj.FotoUtama,&obj.FotoSatu,&obj.FotoDua,&obj.FotoTiga,&obj.FotoEmpat,&obj.IdStatus,&obj.CreatedAt,&obj.UpdatedAt,&obj.DeletedAt)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}

func (r *FashionRepositoryImpl) FindAllByParameter(param map[string]string)(*model.FashionSearchs,error){
	querys := `select distinct b.id_status,b.id_busana,k.nama as city,IFNULL(rs.rating_total,0) as rating ,b.nama_busana as name,b.url_foto_utama as picture,b.harga as price from tb_busana b
	INNER JOIN tb_user u ON b.id_user = u.id_user
	INNER JOIN tb_alamat a ON u.id_user = a.id_user
	LEFT JOIN tb_pengiriman_busana pb ON pb.id_busana = b.id_busana
    LEFT JOIN tb_rating_summary rs ON rs.id_user = u.id_user
	LEFT JOIN tb_kabupaten k ON a.kabupaten = k.id_kabupaten
    WHERE 1=1 `

    querysTrx := ` select t1.* from temp t1 WHERE t1.id_status = 'STS01' UNION
select t2.* from temp2 t2 INNER JOIN tb_transaksi_detail td ON t2.id_busana = td.id_busana
INNER JOIN tb_transaksi t ON t.id_transaksi = td.id_transaksi WHERE 1=1 `

    if(len(param) > 0){
		if(len(param["idUser"]) > 0){
			cond := fmt.Sprintf(`AND b.id_user != '%s'`,param["idUser"])
			querys = fmt.Sprintf("%s%s",querys,cond)
		}
		if(len(param["name"]) > 0){
			cond := `AND b.nama_busana LIKE '%`+param["name"]+`%' `
			querys = fmt.Sprintf("%s%s",querys,cond)
		}
		if(len(param["idCategory"]) > 0){
			cond := fmt.Sprintf(`AND b.id_kategori = '%s' `,param["idCategory"])
			querys = fmt.Sprintf("%s%s",querys,cond)
		}
		if(len(param["location"]) > 0){
			cond := fmt.Sprintf(`AND a.kabupaten = '%s' `,param["location"])
			querys = fmt.Sprintf("%s%s",querys,cond)
		}
		if(len(param["idShippment"]) > 0){
			cond := fmt.Sprintf(`AND pb.id_pengiriman = '%s' `,param["idShippment"])
			querys = fmt.Sprintf("%s%s",querys,cond)
		}
		if(len(param["startDate"]) > 0 && len(param["endDate"]) > 0){
			startDate := param["startDate"]
			endDate := param["endDate"]
			cond := fmt.Sprintf(`AND t.start_date - 86400 not between %s and %s and t.end_date + 86400 not between %s and %s;`,startDate,endDate,startDate,endDate)
			querysTrx = fmt.Sprintf("%s%s",querysTrx,cond)

		}
	}

	queryTotal := fmt.Sprintf(`CREATE TEMPORARY TABLE temp (%s);
CREATE TEMPORARY TABLE temp2 (SELECT * FROM temp);
%s
DROP TEMPORARY TABLE temp;DROP TEMPORARY TABLE temp2;`,querys,querysTrx)

	var objs model.FashionSearchs

	rows, err := r.db.Query(queryTotal)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.FashionSearch

		err = rows.Scan(&obj.IdStatus,&obj.Id,&obj.Location,&obj.Rating,&obj.Name,&obj.UrlImage,&obj.Price)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil

}

func (r *FashionRepositoryImpl) FindByIdUser(idUser string) (*model.FashionSearchs,error) {
	querys := `select distinct b.id_status,b.id_busana,k.nama as city,IFNULL(rs.rating_total,0) as rating ,b.nama_busana as name,b.url_foto_utama as picture,b.harga as price from tb_busana b
INNER JOIN tb_user u ON b.id_user = u.id_user
INNER JOIN tb_alamat a ON u.id_user = a.id_user
LEFT JOIN tb_pengiriman_busana pb ON pb.id_busana = b.id_busana
LEFT JOIN tb_rating_summary rs ON rs.id_user = u.id_user
LEFT JOIN tb_kabupaten k ON a.kabupaten = k.id_kabupaten
WHERE b.id_user = '`+idUser+`'`

	var objs model.FashionSearchs

	rows, err := r.db.Query(querys)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.FashionSearch

		err = rows.Scan(&obj.IdStatus,&obj.Id,&obj.Location,&obj.Rating,&obj.Name,&obj.UrlImage,&obj.Price)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}
