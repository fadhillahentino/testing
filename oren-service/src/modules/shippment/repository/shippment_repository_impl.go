package repository

import (
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/modules/shippment/model"
	"fmt"
	"strings"
)

type ShippmentRepositoryImpl struct {
	db *sql.DB
}

func NewShippmentRepositoryImpl(db *sql.DB) *ShippmentRepositoryImpl {
	return &ShippmentRepositoryImpl{db}
}

func (r *ShippmentRepositoryImpl) SaveShippmentFashion(obj []string, idBusana string) error{
	query := fmt.Sprintf(`DELETE FROM tb_pengiriman_busana WHERE id_busana = '%s' ;INSERT INTO tb_pengiriman_busana (id_busana, id_pengiriman) VALUES `,idBusana)

	for _, v := range obj{
		query = fmt.Sprintf(`%s ('%s', '%s'),`,query,idBusana,v)
	}

	querys :=  query[0:len(query)-1]
	_, err := r.db.Exec(querys)
	if err != nil {
		return err
	}

	defer r.db.Close()

	return nil
}

func (r *ShippmentRepositoryImpl) FindByIdBusana(idBusana string) (string,error){
	query := `SELECT DISTINCT p.pengiriman FROM tb_pengiriman_busana pb
LEFT JOIN tb_pengiriman p ON pb.id_pengiriman = p.id_pengiriman
WHERE pb.id_busana = '%s'`
	querys := fmt.Sprintf(query,idBusana)
	rows, err := r.db.Query(querys)
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}

	defer rows.Close()

	objs := ""
	for rows.Next(){
		var obj string

		err = rows.Scan(&obj)
		if err != nil {
			fmt.Println(err.Error())
			return "",err
		}
		objs = fmt.Sprintf(`%s,%s`,objs,obj)
	}

	objs = strings.Replace(objs,",","",1)

	return objs,nil
}

func (r *ShippmentRepositoryImpl) FindAll() (*model.Shippments,error){
	fmt.Println("repo.shippment.findAll")
	query := `SELECT * FROM tb_pengiriman`

	var objs model.Shippments

	rows, err := r.db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}

	defer rows.Close()

	for rows.Next(){
		var obj model.Shippment

		err = rows.Scan(&obj.IdPengiriman,&obj.Pengiriman)
		if err != nil {
			fmt.Println(err.Error())
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}
func (r *ShippmentRepositoryImpl) FindAllByIdBusana(idBusana string) (*model.Shippments,error){
	query := `SELECT DISTINCT p.id_pengiriman, p.pengiriman FROM tb_pengiriman_busana pb
LEFT JOIN tb_pengiriman p ON pb.id_pengiriman = p.id_pengiriman
WHERE pb.id_busana = '%s'`
	querys := fmt.Sprintf(query,idBusana)
	rows, err := r.db.Query(querys)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}

	defer rows.Close()

	var objs model.Shippments
	for rows.Next(){
		var obj model.Shippment

		err = rows.Scan(&obj.IdPengiriman,&obj.Pengiriman)
		if err != nil {
			fmt.Println(err.Error())
			return nil,err
		}
		objs = append(objs,obj)
	}

	return &objs,nil
}