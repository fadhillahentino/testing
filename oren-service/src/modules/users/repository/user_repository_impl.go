package repository

import (
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/modules/users/model"
	"github.com/fadhilfcr/oren-service/src/util/password"
	"errors"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db}
}

func (r *UserRepositoryImpl) Save(user *model.User) error{
	query := `INSERT INTO tb_user (id_user,nama,no_hp,email,password,foto_user,uid,created_time,updated_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(user.IdUser, user.Name, user.Phone, user.Email, user.Password, user.Foto,user.Uid, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) Update(id string,user *model.User) error{
	query := `UPDATE tb_user SET nama = ?, no_hp = ?, email = ?, password = ?, foto = ?, updated_time = ?) where id_user = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(user.Name, user.Phone,user.Email, user.Password, user.Foto, user.UpdatedAt, user.IdUser)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) Delete(id string) error{
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

func (r *UserRepositoryImpl) FindById(id string) (*model.User,error){
	query := `SELECT id_user, no_hp, nama, email, foto_user, IFNULL(url_foto_user,'') as url_foto_user, password, uid, created_time, updated_time FROM tb_user WHERE id_user = ?`

	var user model.User

	statement, err := r.db.Prepare(query)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	err = statement.QueryRow(id).Scan(&user.IdUser,&user.Phone,&user.Name,&user.Email,&user.Foto,&user.UrlFoto,&user.Password,&user.Uid,&user.CreatedAt,&user.UpdatedAt)
	if err != nil {
		return nil,err
	}

	return &user,nil
}

func (r *UserRepositoryImpl) FindAll() (*model.Users,error){
	query := `SELECT * FROM tb_user`

	var users model.Users

	rows, err := r.db.Query(query)
	if err != nil {
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var user model.User

		err = rows.Scan(&user.IdUser,&user.Name,&user.Phone,&user.Email,&user.Password,&user.Foto,&user.CreatedAt,&user.UpdatedAt)
		if err != nil {
			return nil,err
		}
		users = append(users,user)
	}

	return &users,nil
}

func (r *UserRepositoryImpl) CheckLogin(phone string, inPwd string) (string,error){
	var dbPwd string
	var iduser string

	query := `SELECT id_user, password FROM tb_user where no_hp = ? LIMIT 1`
	statement, err := r.db.Prepare(query)
	if err != nil {
		return "",err
	}

	defer statement.Close()

	err = statement.QueryRow(phone).Scan(&iduser,&dbPwd)
	if err != nil {
		return "",err
	}

	result := password.CheckPasswordHash(inPwd,dbPwd)
	if !result {
		return "",errors.New("Not Match")
	}
	return iduser,nil
}

func (r *UserRepositoryImpl) CheckRegistration(phone string, email string) bool{
	var count int

	query := `SELECT count(1) FROM tb_user where no_hp = ? OR email = ? LIMIT 1`
	statement, err := r.db.Prepare(query)
	if err != nil {
		return false
	}

	defer statement.Close()

	err = statement.QueryRow(phone,email).Scan(&count)
	if err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *UserRepositoryImpl) GetAuthFirebase(idUser string) (*model.User,error){
	query := `select nama,no_hp,email,password,uid,url_foto_user from tb_user WHERE id_user = ?;`
	statement, err := r.db.Prepare(query)
	if err != nil {
		return nil,err
	}

	defer statement.Close()

	var obj *model.User

	err = statement.QueryRow(idUser).Scan(&obj.Name,&obj.Phone,&obj.Email,&obj.Password,&obj.Uid,&obj.UrlFoto)
	if err != nil {
		return nil,err
	}

	return obj,nil
}







