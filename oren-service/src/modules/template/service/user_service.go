package service

import (
	"net/http"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"database/sql"
	"github.com/fadhilfcr/oren-service/src/util/password"

	genericRepo "github.com/fadhilfcr/oren-service/src/modules/base/repository"

	stringUtil "github.com/fadhilfcr/oren-service/src/util/string"
	constant "github.com/fadhilfcr/oren-service/src/util"
	userRepo "github.com/fadhilfcr/oren-service/src/modules/users/repository"
	"github.com/fadhilfcr/oren-service/src/modules/users/model"

	addrModel "github.com/fadhilfcr/oren-service/src/modules/address/model"
	addrRepo "github.com/fadhilfcr/oren-service/src/modules/address/repository"
)

func Login_service(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusNotFound
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		phone,_ := jsonparser.GetString(body,"phone") //089123123123
		pwd,_ := jsonparser.GetString(body,"pwd") //fadhil123

		userRepoImpl := userRepo.NewUserRepositoryImpl(db)

		err = userRepoImpl.CheckLogin(phone,pwd)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		code = http.StatusOK
		w.WriteHeader(code)
	}
}

func Registration_service(db *sql.DB)func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		name,_ := jsonparser.GetString(body,"name")
		phone,_ := jsonparser.GetString(body,"phone") //089123123123
		email,_ := jsonparser.GetString(body,"email")
		pwd,_ := jsonparser.GetString(body,"pwd") //fadhil123
		photo,_ := jsonparser.GetString(body,"photo")

		address,_ := jsonparser.GetString(body,"address")
		provinsi,_ := jsonparser.GetString(body,"provinsi")
		kabupaten,_ := jsonparser.GetString(body,"kabupaten")
		kecamatan,_ := jsonparser.GetString(body,"kecamatan")
		kodePos,_ := jsonparser.GetString(body,"kodePos")

		userRepoImpl := userRepo.NewUserRepositoryImpl(db)
		addrRepoImpl := addrRepo.NewAddressRepositoryImpl(db)
		genericRepoImpl := genericRepo.NewGenericRepositoryImpl(db)

		check := userRepoImpl.CheckRegistration(phone,email)
		if !check{
			w.WriteHeader(http.StatusConflict)
			return
		}

		maxId,_ := genericRepoImpl.FindMaxId(constant.TAG_USER_PK,constant.TAG_USER_TBL,constant.TAG_USER)
		maxIdAlamat,_ := genericRepoImpl.FindMaxId(constant.TAG_ADDRESS_PK,constant.TAG_ADDRESS_TBL,constant.TAG_ADDRESS)

		idUser,err := stringUtil.TableIdFormatter(constant.TAG_USER,maxId)
		idAlamat,err := stringUtil.TableIdFormatter(constant.TAG_ADDRESS,maxIdAlamat)

		user := model.NewUser()
		user.IdUser = idUser
		user.Nama = name
		user.NoHp = phone
		user.Email = email
		user.Password,_ = password.HashPassword(pwd)
		user.Foto = photo

		addr := addrModel.NewAddress()
		addr.IdUser = idUser
		addr.IdAlamat = idAlamat
		addr.Alamat = address
		addr.Provinsi = provinsi
		addr.Kabupaten = kabupaten
		addr.Kecamatan = kecamatan
		addr.KodePos = kodePos

		err = userRepoImpl.Save(user)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		err = addrRepoImpl.Save(addr)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
	}
}