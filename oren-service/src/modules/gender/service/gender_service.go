package service

import (
	"net/http"
	"database/sql"
	repo "github.com/fadhilfcr/oren-service/src/modules/gender/repository"
	"encoding/json"
)

func FindAll(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusNotFound

		userRepoImpl := repo.NewGenderRepositoryImpl(db)

		objs,err := userRepoImpl.FindAll()
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