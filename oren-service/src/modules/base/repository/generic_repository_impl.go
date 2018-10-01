package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"strconv"
)

type GenericRepositoryImpl struct {
	db *sql.DB
}

func NewGenericRepositoryImpl(db *sql.DB) *GenericRepositoryImpl {
	return &GenericRepositoryImpl{db}
}

func (r *GenericRepositoryImpl) FindMaxId(primaryField string,tblName string,tblTag string) (int,error){
	query := fmt.Sprintf(`SELECT IFNULL(MAX(%s),"0") FROM %s`,primaryField,tblName)

	var strMaxId string

	statement, err := r.db.Prepare(query)
	if err != nil {
		return 0,err
	}

	defer statement.Close()

	err = statement.QueryRow().Scan(&strMaxId)
	if err != nil {
		return 0,err
	}

	strData := strings.Replace(strMaxId,tblTag,"",-1)
	maxId,_ := strconv.Atoi(strData)

	maxId++

	return maxId,nil
}
