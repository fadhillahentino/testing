package repository

import (
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/modules/rating/model"
	"fmt"
	"time"
)

type RatingRepositoryImpl struct {
	db *sql.DB
}

func NewRatingRepositoryImpl(db *sql.DB) *RatingRepositoryImpl {
	return &RatingRepositoryImpl{db}
}

func (r *RatingRepositoryImpl) Save(obj *model.Rating) error{
	query := `INSERT INTO tb_rating (id_rating, id_transaksi, id_user, ulasan, respon_chat, fleksibel, tepat_waktu, perawatan_busana, friendly, total, flag, created_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.IdRating,obj.IdTransaksi,obj.IdUser,obj.Desc,obj.RespChat,obj.Flexible,obj.TepatWaktu,obj.Perawatan,obj.Friendly,obj.Total,obj.Flag,time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func (r *RatingRepositoryImpl) SaveSum(obj *model.RatingSummary) error{
	qDel := `DELETE FROM tb_rating_summary WHERE id_user = '`+obj.IdUser+`' ;`
	_, err := r.db.Exec(qDel)
	if err != nil {
		return err
	}
	query := `INSERT INTO tb_rating_summary (id_user, respon_chat, fleksibel, tepat_waktu, perawatan_busana, friendly, jml_data, total, flag, respon_chat_rat, fleksibel_rat, tepat_waktu_rat, perawatan_busana_rat, friendly_rat, rating_total, created_time, updated_time) VALUES (?, ?, ?, ?, ?,?, ?, ?, ?, ?,?, ?, ?, ?, ?,?,?);`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(obj.IdUser,obj.RespChat,obj.Flexible,obj.TepatWaktu,obj.Perawatan,obj.Friendly,obj.JmlData,obj.Total,0,obj.RespChatRat,obj.FlexibleRat,obj.TepatWaktuRat,obj.PerawatanRat,obj.FriendlyRat,obj.TotalRat,time.Now().Unix(),time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func (r *RatingRepositoryImpl) FindByIdTrx(idTrx string,flagTmp int64) (string,int64,error){
	id := `id_penyewa`
	var flag int64
	flag = 1
	if(flagTmp == 1){
		id = `id_pemilik`
		flag = 0
	}
	query := fmt.Sprintf(`SELECT %s FROM tb_transaksi WHERE id_transaksi = ?`,id)

	var obj model.Rating

	statement, err := r.db.Prepare(query)
	if err != nil {
		return "",0,err
	}

	defer statement.Close()

	err = statement.QueryRow(idTrx).Scan(&obj.IdUser)
	if err != nil {
		return "",0,err
	}

	return obj.IdUser,flag,nil
}

/*func (r *UserRepositoryImpl) FindAll() (*model.Users,error){
	query := `SELECT * FROM tb_user`

	var users model.Users

	rows, err := r.db.Query(query)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var user model.User

		err = rows.Scan(&user.IdUser,&user.Nama,&user.NoHp,&user.Email,&user.Password,&user.Foto,&user.CreatedAt,&user.UpdatedAt)
		if err != nil {
			return nil,err
		}
		users = append(users,user)
	}

	return &users,nil
}*/

func (r *RatingRepositoryImpl) FindSummaryByIdUser(idUser string) (*model.RatingSummary,error){
	query := `SELECT id_user, respon_chat, fleksibel, tepat_waktu, perawatan_busana, friendly, total,
respon_chat_rat, fleksibel_rat, tepat_waktu_rat, perawatan_busana_rat, friendly_rat, rating_total,jml_data
FROM tb_rating_summary WHERE id_user = ?`

	var obj model.RatingSummary

	statement, err := r.db.Prepare(query)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow(idUser).Scan(&obj.IdUser,&obj.RespChat,&obj.Flexible,&obj.TepatWaktu,&obj.Perawatan,&obj.Friendly,&obj.Total,
		&obj.RespChatRat,&obj.FlexibleRat,&obj.TepatWaktuRat,&obj.PerawatanRat,&obj.FriendlyRat,&obj.TotalRat,&obj.JmlData)
	if err != nil {
		return nil,err
	}

	return &obj,nil
}

func (r *RatingRepositoryImpl) FindAccountByIdUser(idUser string) (*model.Rating,error){
	query := `select u.id_user,u.nama,rs.respon_chat_rat,rs.fleksibel_rat	,rs.tepat_waktu_rat	,rs.perawatan_busana_rat	,rs.friendly_rat
,rs.rating_total,COUNT(r.id_user)  from tb_user u LEFT JOIN tb_rating_summary rs ON u.id_user = rs.id_user
LEFT JOIN tb_rating r ON u.id_user = r.id_user WHERE u.id_user = ?
GROUP BY u.id_user,u.nama,rs.respon_chat_rat,rs.fleksibel_rat	,rs.tepat_waktu_rat	,rs.perawatan_busana_rat	,rs.friendly_rat	,rs.rating_total;`
	var obj model.Rating

	statement, err := r.db.Prepare(query)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow(idUser).Scan(&obj.IdUser,&obj.Nama,&obj.RespChat,&obj.Flexible,&obj.TepatWaktu,&obj.Perawatan,&obj.Friendly,&obj.Total,&obj.JmlData)
	if err != nil {
		return nil,err
	}

	return &obj,nil
}

func (r *RatingRepositoryImpl) FindComments(idUser string) (*model.Ratings,error){
	querys := `select ulasan,case when flag = 1 then 'Penyewa' else 'Pemilik' end as status from tb_rating WHERE id_user = '`+idUser+`' order by created_time desc;`

	var objs model.Ratings

	rows, err := r.db.Query(querys)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.Rating

		err = rows.Scan(&obj.Desc,&obj.Status)
		if err != nil {
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}

