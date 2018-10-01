package repository

import (
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/modules/transaction/model"
	"log"
	"time"
)

type TrxRepositoryImpl struct {
	db *sql.DB
}

func NewTrxRepositoryImpl(db *sql.DB) *TrxRepositoryImpl {
	return &TrxRepositoryImpl{db}
}

func (r *TrxRepositoryImpl) Save(obj *model.Transaction) error{
	query := `INSERT INTO tb_transaksi(id_transaksi, id_penyewa, id_pemilik, start_date, end_date, total_harga,
bukti_transfer, bukti_pengiriman, bukti_pengiriman_kembali, url_bukti_transfer, url_bukti_pengiriman,
url_bukti_pengiriman_kembali, id_pengiriman, id_status, id_alamat, created_time, updated_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.IdTransaksi, obj.IdPenyewa, obj.IdPemilik, obj.StartDate,obj.EndDate,obj.TotalHarga,
		obj.BuktiTransfer,obj.BuktiPengiriman,obj.BuktiPengirimanKembali,obj.UrlBuktiTransfer,obj.UrlBuktiPengiriman,
			obj.UrlBuktiPengirimanKembali,obj.IdPengiriman,obj.IdStatus, obj.IdAlamat, obj.CreatedAt, obj.UpdatedAt)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func (r *TrxRepositoryImpl) Update(id string,obj *model.Transaction) error{
	query := `UPDATE tb_transaksi SET id_penyewa=?,id_pemilik=?,start_date=?,end_date=?,total_harga=?,bukti_transfer=?,
bukti_pengiriman=?,bukti_pengiriman_kembali=?,url_bukti_transfer=?,url_bukti_pengiriman=?,url_bukti_pengiriman_kembali=?,
id_pengiriman=?,id_status=?, id_alamat=?, updated_time=? where id_user = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.IdPenyewa,obj.IdPemilik,obj.StartDate,obj.EndDate,obj.TotalHarga,obj.BuktiTransfer,
		obj.BuktiPengiriman,obj.BuktiPengirimanKembali,obj.UrlBuktiTransfer,obj.UrlBuktiPengiriman,obj.UrlBuktiPengirimanKembali,
		obj.IdPengiriman,obj.IdStatus,obj.IdAlamat,obj.UpdatedAt, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *TrxRepositoryImpl) UpdateStatus(id string,idStatus string) error{
	query := `UPDATE tb_transaksi SET id_status=?, updated_time=? where id_transaksi = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(idStatus,time.Now().Unix(), id)
	if err != nil {
		return err
	}

	return nil
}

func (r *TrxRepositoryImpl) UpdateStatusApprove(id string,idStatus string,idRekening string, biayaPengiriman float64,idAlamatPemilik string) error{
	query := `UPDATE tb_transaksi SET total_harga = total_harga+? ,id_status=?, harga_pengiriman = ?, id_rekening=?, updated_time=?,id_alamat_pemilik = ? where id_transaksi = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(biayaPengiriman,idStatus,biayaPengiriman,idRekening,time.Now().Unix(),idAlamatPemilik, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *TrxRepositoryImpl) Delete(id string) error{
	query := `DELETE tb_user where id_user = ?`

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

func (r *TrxRepositoryImpl) FindById(id string) (*model.Transaction,error){
	query := `SELECT id_transaksi, id_penyewa, id_pemilik, start_date, end_date, total_harga, IFNULL(bukti_transfer,''), IFNULL(bukti_pengiriman,''), IFNULL(bukti_pengiriman_kembali,''), IFNULL(url_bukti_transfer,''), IFNULL(nama_bank_transfer,''), IFNULL(atas_nama_transfer,''), nominal_transfer, IFNULL(url_bukti_pengiriman,''), IFNULL(url_bukti_pengiriman_kembali,''), IFNULL(id_pengiriman,''), IFNULL(id_pengiriman_kembali,''), IFNULL(id_status,''), IFNULL(id_alamat,''), IFNULL(id_alamat_pemilik,''), created_time, updated_time, harga_pengiriman, IFNULL(id_rekening,'') FROM tb_transaksi  WHERE id_transaksi = ?`

	var obj model.Transaction

	statement, err := r.db.Prepare(query)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow(id).Scan(&obj.IdTransaksi, &obj.IdPenyewa, &obj.IdPemilik, &obj.StartDate,&obj.EndDate,&obj.TotalHarga,
		&obj.BuktiTransfer,&obj.BuktiPengiriman,&obj.BuktiPengirimanKembali,&obj.UrlBuktiTransfer,&obj.NamaBank,&obj.AtasNama,&obj.NominaTransfer,
			&obj.UrlBuktiPengiriman,&obj.UrlBuktiPengirimanKembali,&obj.IdPengiriman,&obj.IdPengirimanKembali,&obj.IdStatus, &obj.IdAlamat,
		&obj.IdAlamatPemilik,&obj.CreatedAt, &obj.UpdatedAt,&obj.HargaPengiriman,&obj.IdRekening)
	if err != nil {
		return nil,err
	}

	return &obj,nil
}

func (r *TrxRepositoryImpl) FindAll() (*model.Transactions,error){
	query := `SELECT * FROM tb_transaksi`

	var objs model.Transactions

	rows, err := r.db.Query(query)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.Transaction

		err = rows.Scan(&obj.IdTransaksi, &obj.IdPenyewa, &obj.IdPemilik, &obj.StartDate,&obj.EndDate,&obj.TotalHarga,
			&obj.BuktiTransfer,&obj.BuktiPengiriman,&obj.BuktiPengirimanKembali,&obj.UrlBuktiTransfer,&obj.UrlBuktiPengiriman,
			&obj.UrlBuktiPengirimanKembali,&obj.IdPengiriman,&obj.IdStatus, &obj.IdAlamat, &obj.CreatedAt, &obj.UpdatedAt)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}

func (r *TrxRepositoryImpl) SaveDetail(obj *model.Transaction) error{
	query := `INSERT INTO tb_transaksi_detail(id_transaksi,id_busana) VALUES (?, ?)`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.IdTransaksi, obj.IdBusana)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func (r *TrxRepositoryImpl) ListNotificationUser(idUser string)(*model.Notifications,error){
	/*querys := `select t.id_transaksi,b.nama_busana,t.id_status,u.nama as penyewa,up.nama as pemilik,
DATE_FORMAT(from_unixtime(t.start_date), '%d/%m/%Y') as start_date,
DATE_FORMAT(from_unixtime(t.end_date), '%d/%m/%Y') as end_date from tb_transaksi t
LEFT JOIN tb_transaksi_detail td
ON t.id_transaksi = td.id_transaksi
LEFT JOIN tb_busana b
ON b.id_busana = td.id_busana
LEFT JOIN tb_user u
ON u.id_user = t.id_penyewa
LEFT JOIN tb_user up
ON up.id_user = t.id_pemilik
WHERE (t.id_status IN ('STS11','STS15','STS18','STS19','STS20','STS21','STS23')
AND t.id_pemilik = '%s') OR  (t.id_status IN ('STS12','STS13','STS14','STS16','STS17','STS20','STS22','STS23')
AND t.id_penyewa = '%s') ORDER BY t.updated_time desc;`*/
querys := `select * from (select t.id_transaksi,b.nama_busana,t.id_status,s.status,u.nama as penyewa,up.nama as pemilik,DATE_FORMAT(from_unixtime(t.start_date), '%d/%m/%Y') as start_date,DATE_FORMAT(from_unixtime(t.end_date), '%d/%m/%Y') as end_date,t.updated_time,0 as flag from tb_transaksi t
LEFT JOIN tb_transaksi_detail td
ON t.id_transaksi = td.id_transaksi
LEFT JOIN tb_busana b
ON b.id_busana = td.id_busana
LEFT JOIN tb_user u
ON u.id_user = t.id_penyewa
LEFT JOIN tb_user up
ON up.id_user = t.id_pemilik
LEFT JOIN tb_status s
ON s.id_status = t.id_status
WHERE t.id_pemilik = '`+idUser+`'
UNION
select t.id_transaksi,b.nama_busana,t.id_status,s.status,u.nama as penyewa,up.nama as pemilik,DATE_FORMAT(from_unixtime(t.start_date), '%d/%m/%Y') as start_date,DATE_FORMAT(from_unixtime(t.end_date), '%d/%m/%Y') as end_date,t.updated_time,1 from tb_transaksi t
LEFT JOIN tb_transaksi_detail td
ON t.id_transaksi = td.id_transaksi
LEFT JOIN tb_busana b
ON b.id_busana = td.id_busana
LEFT JOIN tb_user u
ON u.id_user = t.id_penyewa
LEFT JOIN tb_user up
ON up.id_user = t.id_pemilik
LEFT JOIN tb_status s
ON s.id_status = t.id_status
WHERE t.id_penyewa = '`+idUser+`') a ORDER BY a.updated_time desc;`

	var objs model.Notifications

	rows, err := r.db.Query(querys)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.Notification

		err = rows.Scan(&obj.IdTransaksi, &obj.NamaBusana, &obj.IdStatus, &obj.Status, &obj.Penyewa, &obj.Pemilik, &obj.StartDate, &obj.EndDate,&obj.UpdatedAt,&obj.Flag)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}

func (r *TrxRepositoryImpl) RentDetail(idTrx string)(*model.RentDetail,error){
	querys := `select up.id_user,t.id_transaksi,b.nama_busana,t.id_status,u.nama as penyewa,up.nama as pemilik,
DATE_FORMAT(from_unixtime(t.start_date), '%d/%m/%Y') as start_date,
DATE_FORMAT(from_unixtime(t.end_date), '%d/%m/%Y') as end_date,
IFNULL(rsu.rating_total,0) as rating_penyewa,IFNULL(rsp.rating_total,0) as rating_pemilik,
IFNULL(u.url_foto_user,'') as url_foto_user_penyewa,IFNULL(up.url_foto_user,'') as url_foto_user_pemilik,
b.url_foto_utama,pe.pengiriman,IFNULL(t.harga_pengiriman,0) as harga_pengiriman,b.harga,
IFNULL(b.deposit,0) as deposit,IFNULL(t.total_harga,0)-IFNULL(b.deposit,0) as harga_sewa,
IFNULL(t.total_harga,0) as total_harga,
a.alamat_lengkap,a.kecamatan,k.nama as kabupaten,p.nama as provinsi,a.kode_pos,
u.no_hp from tb_transaksi t
LEFT JOIN tb_transaksi_detail td
ON t.id_transaksi = td.id_transaksi
LEFT JOIN tb_busana b
ON b.id_busana = td.id_busana
LEFT JOIN tb_user u
ON u.id_user = t.id_penyewa
LEFT JOIN tb_user up
ON up.id_user = t.id_pemilik
LEFT JOIN tb_rating_summary rsu
ON rsu.id_user = u.id_user
LEFT JOIN tb_rating_summary rsp
ON rsp.id_user = up.id_user
LEFT JOIN tb_pengiriman pe
ON pe.id_pengiriman = t.id_pengiriman
LEFT JOIN tb_alamat a
ON a.id_alamat = t.id_alamat
LEFT JOIN tb_provinsi p
ON a.provinsi = p.id_provinsi
LEFT JOIN tb_kabupaten k
ON a.kabupaten = k.id_kabupaten
WHERE t.id_transaksi = '`+idTrx+`';`
	var obj model.RentDetail

	statement, err := r.db.Prepare(querys)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow().Scan(&obj.IdUserPemilik,&obj.IdUserPemilik,&obj.IdTransaksi, &obj.NamaBusana, &obj.IdStatus, &obj.Penyewa, &obj.Pemilik, &obj.StartDate,
		&obj.EndDate, &obj.RatingPenyewa, &obj.RatingPemilik, &obj.UrlPenyewa, &obj.UrlPemilik, &obj.UrlBusana, &obj.Pengiriman,
			&obj.HargaPengiriman,&obj.Harga, &obj.Deposit, &obj.HargaSewa, &obj.TotalHarga, &obj.Alamat, &obj.Kecamatan,&obj.Kabupaten,&obj.Provinsi,&obj.KodePos, &obj.NoHp)
	if err != nil {
		return nil,err
	}

	return &obj,nil
}

func (r *TrxRepositoryImpl) ConfirmBilling(idTrx string)(*model.RentDetail,error){
	querys := `select t.id_transaksi,b.nama_busana,t.id_status,u.nama as penyewa,up.nama as pemilik,
DATE_FORMAT(from_unixtime(t.start_date), '%d/%m/%Y') as start_date,
DATE_FORMAT(from_unixtime(t.end_date), '%d/%m/%Y') as end_date,IFNULL(t.harga_pengiriman,0) as harga_pengiriman,b.harga,
IFNULL(b.deposit,0) as deposit,IFNULL(t.total_harga,0)-IFNULL(b.deposit,0)-IFNULL(t.harga_pengiriman,0) as harga_sewa,IFNULL(t.total_harga,0) as total_harga,
r.nomor_rekening,r.nama_bank,r.atas_nama,pe.pengiriman from tb_transaksi t
LEFT JOIN tb_transaksi_detail td
ON t.id_transaksi = td.id_transaksi
LEFT JOIN tb_busana b
ON b.id_busana = td.id_busana
LEFT JOIN tb_user u
ON u.id_user = t.id_penyewa
LEFT JOIN tb_user up
ON up.id_user = t.id_pemilik
LEFT JOIN tb_rekening r
ON r.id_rekening = t.id_rekening
LEFT JOIN tb_pengiriman pe
ON pe.id_pengiriman = t.id_pengiriman
WHERE t.id_transaksi = '`+idTrx+`';`

	var obj model.RentDetail

	statement, err := r.db.Prepare(querys)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow().Scan(&obj.IdTransaksi, &obj.NamaBusana, &obj.IdStatus, &obj.Penyewa, &obj.Pemilik, &obj.StartDate,
		&obj.EndDate, &obj.HargaPengiriman,&obj.Harga, &obj.Deposit, &obj.HargaSewa, &obj.TotalHarga, &obj.NomorRekening, &obj.NamaBank, &obj.AtasNama,&obj.Pengiriman )
	if err != nil {
		return nil,err
	}

	return &obj,nil
}

func (r *TrxRepositoryImpl) ConfirmBillingSave(obj *model.Transaction) error{
	query := `UPDATE tb_transaksi SET bukti_transfer = ?, url_bukti_transfer = ?, nama_bank_transfer = ?, atas_nama_transfer = ?, nominal_transfer = ?,id_status=?, updated_time=? where id_transaksi = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.BuktiTransfer,obj.UrlBuktiTransfer,obj.NamaBank,obj.AtasNama,obj.NominaTransfer,obj.IdStatus ,time.Now().Unix(),obj.IdTransaksi)
	if err != nil {
		return err
	}

	return nil
}

func (r *TrxRepositoryImpl) SendProductSave(obj *model.Transaction) error{
	query := `UPDATE tb_transaksi SET bukti_pengiriman = ?, url_bukti_pengiriman = ?, id_status=?, updated_time=? where id_transaksi = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.Resi,obj.UrlBuktiPengiriman,obj.IdStatus ,time.Now().Unix(),obj.IdTransaksi)
	if err != nil {
		return err
	}

	return nil
}

func (r *TrxRepositoryImpl) ReceiptPayment(idTrx string)(*model.RentDetail,error){
	querys := `SELECT 	DATE_FORMAT(from_unixtime(start_date-86400), '%d/%m/%Y') as start_date,
nama_bank_transfer,
atas_nama_transfer,
nominal_transfer,
url_bukti_transfer
FROM tb_transaksi WHERE id_transaksi = '`+idTrx+`';`

	var obj model.RentDetail

	statement, err := r.db.Prepare(querys)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow().Scan(&obj.StartDate,
		&obj.NamaBank, &obj.AtasNama,&obj.NominaTransfer,&obj.UrlBuktiTransfer )
	if err != nil {
		return nil,err
	}

	return &obj,nil
}

func (r *TrxRepositoryImpl) SendProduct(idTrx string)(*model.RentDetail,error){
	querys := `SELECT	DATE_FORMAT(from_unixtime(t.end_date+86400), '%d/%m/%Y') as end_date,p.pengiriman
FROM tb_transaksi t LEFT JOIN tb_pengiriman p ON t.id_pengiriman = p.id_pengiriman
WHERE id_transaksi = '`+idTrx+`';`

	var obj model.RentDetail

	statement, err := r.db.Prepare(querys)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow().Scan(&obj.EndDate,
		&obj.Pengiriman )
	if err != nil {
		return nil,err
	}

	return &obj,nil
}

func (r *TrxRepositoryImpl) ReceiveProduct(idTrx string)(*model.RentDetail,error){
	querys := `SELECT	DATE_FORMAT(from_unixtime(t.end_date+86400), '%d/%m/%Y') as end_date,p.pengiriman,t.url_bukti_pengiriman,t.bukti_pengiriman
FROM tb_transaksi t LEFT JOIN tb_pengiriman p ON t.id_pengiriman = p.id_pengiriman
WHERE id_transaksi = '`+idTrx+`';`

	var obj model.RentDetail

	statement, err := r.db.Prepare(querys)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow().Scan(&obj.EndDate,
		&obj.Pengiriman,&obj.UrlBuktiPengiriman,&obj.BuktiPengiriman )
	if err != nil {
		return nil,err
	}

	return &obj,nil
}

func (r *TrxRepositoryImpl) RetriveProduct(idTrx string)(*model.RentDetail,error){
	querys := `SELECT CONCAT(a.alamat_lengkap,', ',a.kecamatan,', ',k.nama,' - ',p.nama,', ',a.kode_pos) as alamat FROM tb_transaksi t LEFT JOIN tb_alamat a ON t.id_alamat_pemilik = a.id_alamat
LEFT JOIN tb_provinsi p ON p.id_provinsi = a.provinsi
LEFT JOIN tb_kabupaten k ON k.id_kabupaten = a.kabupaten
WHERE id_transaksi = '`+idTrx+`';`

	var obj model.RentDetail

	statement, err := r.db.Prepare(querys)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow().Scan(&obj.Alamat)
	if err != nil {
		return nil,err
	}

	return &obj,nil
}

func (r *TrxRepositoryImpl) RetriveProductSave(obj *model.Transaction) error{
	query := `UPDATE tb_transaksi SET id_pengiriman_kembali = ?,bukti_pengiriman_kembali = ?, url_bukti_pengiriman_kembali = ?, id_status=?, updated_time=? where id_transaksi = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.IdPengirimanKembali,obj.BuktiPengirimanKembali,obj.UrlBuktiPengirimanKembali,obj.IdStatus ,time.Now().Unix(),obj.IdTransaksi)
	if err != nil {
		return err
	}

	return nil
}

func (r *TrxRepositoryImpl) ReceiveProductOwner(idTrx string)(*model.RentDetail,error){
	querys := `SELECT	DATE_FORMAT(from_unixtime(t.end_date+86400), '%d/%m/%Y') as end_date,p.pengiriman,IFNULL(t.url_bukti_pengiriman_kembali,'') as url,IFNULL(t.bukti_pengiriman_kembali,'') as bukti
FROM tb_transaksi t LEFT JOIN tb_pengiriman p ON t.id_pengiriman_kembali = p.id_pengiriman
WHERE id_transaksi = '`+idTrx+`';`

	var obj model.RentDetail

	statement, err := r.db.Prepare(querys)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow().Scan(&obj.EndDate,
		&obj.Pengiriman,&obj.UrlBuktiPengirimanKembali,&obj.BuktiPengirimanKembali)
	if err != nil {
		return nil,err
	}

	return &obj,nil
}