package service

import (
	"net/http"
	"encoding/json"
	"database/sql"

	repo "github.com/fadhilfcr/oren-service/src/modules/shippment/repository"
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/buger/jsonparser"
)

func FindByIdBusana(db *sql.DB,idBusana string) (string,error){

	repoImpl := repo.NewShippmentRepositoryImpl(db)

	resp,err := repoImpl.FindByIdBusana(idBusana)
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}
	return resp,nil
}

func FindAll(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusNotFound
		fmt.Println("service.shippment.findAll")
		repoImpl := repo.NewShippmentRepositoryImpl(db)

		objs,err := repoImpl.FindAll()
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(code)
			return
		}

		resp,err := json.Marshal(objs)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
		w.Write(resp)
	}
}

func FindAllByIdBusana(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusNotFound

		body,_ := ioutil.ReadAll(r.Body)
		idBusana,_ := jsonparser.GetString(body,"IdBusana")

		fmt.Println("service.shippment.findAll")
		repoImpl := repo.NewShippmentRepositoryImpl(db)

		objs,err := repoImpl.FindAllByIdBusana(idBusana)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(code)
			return
		}

		resp,err := json.Marshal(objs)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
		w.Write(resp)
	}
}

func SaveShippmentFashion(db *sql.DB,idShippment string,idBusana string)error{

	arrStr := strings.Split(idShippment,",")

	repoImpl := repo.NewShippmentRepositoryImpl(db)

	err := repoImpl.SaveShippmentFashion(arrStr,idBusana)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
