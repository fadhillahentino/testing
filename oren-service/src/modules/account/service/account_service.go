package service

import (
	"net/http"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"database/sql"

	accRepo "github.com/fadhilfcr/oren-service/src/modules/account/repository"
	"encoding/json"
)

func FindAccountById(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		body,_ := ioutil.ReadAll(r.Body)

		IdUser,_ := jsonparser.GetString(body,"IdUser")

		code := http.StatusNotFound

		repoImpl := accRepo.NewAccountRepositoryImpl(db)

		objs,err := repoImpl.FindById(IdUser)
		if err != nil {
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