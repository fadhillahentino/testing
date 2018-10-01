package service

import (
	"database/sql"
	"net/http"
	"encoding/json"
	repo "github.com/fadhilfcr/oren-service/src/modules/address/repository"
	"io/ioutil"
	"github.com/buger/jsonparser"
)

func FindAddressByIdUser(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		body,_ := ioutil.ReadAll(r.Body)

		idUser,_ := jsonparser.GetString(body,"IdUser")

		code := http.StatusNotFound

		repoImpl := repo.NewAddressRepositoryImpl(db)

		objs,err := repoImpl.FindByUserId(idUser)
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
