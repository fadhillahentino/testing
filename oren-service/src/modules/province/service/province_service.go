package service

import (
	"net/http"
	"database/sql"
	repo "github.com/fadhilfcr/oren-service/src/modules/province/repository"
	"encoding/json"
)

func FindAllProvince(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusNotFound

		repoImpl := repo.NewProvinceRepositoryImpl(db)

		objs,err := repoImpl.FindAll()
		if err != nil {
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
func FindAllCity(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusNotFound

		repoImpl := repo.NewProvinceRepositoryImpl(db)

		objs,err := repoImpl.FindAllCity()
		if err != nil {
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