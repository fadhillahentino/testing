package service

import (
	"net/http"
	"database/sql"
	repo "github.com/fadhilfcr/oren-service/src/modules/category/repository"
	"encoding/json"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"fmt"
)

func FindbyGender(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusNotFound

		body,_ := ioutil.ReadAll(r.Body)
		idGender,_ := jsonparser.GetString(body,"IdGender")

		userRepoImpl := repo.NewCategoryRepositoryImpl(db)

		objs,err := userRepoImpl.FindbyGender(idGender)
		if err != nil || len(*objs) <= 0 {
			w.WriteHeader(code)
			return
		}

		if len(*objs) <= 0{
			w.WriteHeader(code)
			return
		}

		resp,err := json.Marshal(objs)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		fmt.Println(string(resp))
		code = http.StatusOK
		w.WriteHeader(code)
		w.Write(resp)
	}
}

func FindAll(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusNotFound

		userRepoImpl := repo.NewCategoryRepositoryImpl(db)

		objs,err := userRepoImpl.FindAll()
		if err != nil || len(*objs) <= 0 {
			w.WriteHeader(code)
			return
		}

		if len(*objs) <= 0{
			w.WriteHeader(code)
			return
		}

		resp,err := json.Marshal(objs)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
		w.Write(resp)
	}
}