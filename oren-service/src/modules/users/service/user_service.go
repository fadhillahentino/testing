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
	"fmt"
	"encoding/json"
	"firebase.google.com/go/auth"
	"context"
	"google.golang.org/api/option"
	"firebase.google.com/go"
)

func Login_service(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusNotFound
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		phone,_ := jsonparser.GetString(body,"Phone") //089123123123
		pwd,_ := jsonparser.GetString(body,"Pass") //fadhil123

		userRepoImpl := userRepo.NewUserRepositoryImpl(db)

		idUser,err := userRepoImpl.CheckLogin(phone,pwd)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		user,err := userRepoImpl.FindById(idUser)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		resp,_ := json.Marshal(user)
		fmt.Println(string(resp))
		code = http.StatusOK
		w.WriteHeader(code)
		w.Write(resp)
	}
}

func Registration_service(db *sql.DB,ctx context.Context)func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		name,_ := jsonparser.GetString(body,"Name")
		phone,_ := jsonparser.GetString(body,"Phone") //089123123123
		email,_ := jsonparser.GetString(body,"Email")
		pwd,_ := jsonparser.GetString(body,"Pass") //fadhil123
		photo,_ := jsonparser.GetString(body,"Photo")

		address,_ := jsonparser.GetString(body,"Address")
		provinsi,_ := jsonparser.GetString(body,"Provinsi")
		kabupaten,_ := jsonparser.GetString(body,"Kabupaten")
		kecamatan,_ := jsonparser.GetString(body,"Kecamatan")
		kodePos,_ := jsonparser.GetString(body,"KodePos")

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

		opt := option.WithCredentialsFile("oren-7dc54-firebase-adminsdk-2bnrn-7af4a598e7.json")
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil{
			fmt.Println(err.Error())
			w.WriteHeader(code)
			return
		}
		client,err := app.Auth(ctx)
		if err != nil{
			fmt.Println(err.Error())
			w.WriteHeader(code)
			return
		}
		pass,_ := password.HashPassword(pwd)
		params := (&auth.UserToCreate{}).
			Email(email).
			EmailVerified(false).
			PhoneNumber(phoneFormat(phone)).
			Password(pass).
			DisplayName(name).
			PhotoURL("http://128.199.129.191:9000/trx0002/1527953701_IMG-20180602-WA0004.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=5TKTJZDVIWJVZKDMLAS1%2F20180602%2Fap-southeast-1%2Fs3%2Faws4_request&X-Amz-Date=20180602T153509Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=0a63488e1f869fda52b83f7b362c68796466893eaacf42379783a489863172d6").
			Disabled(false)
		a,err := client.CreateUser(ctx,params)
		if err != nil{
			fmt.Println(err.Error())
			w.WriteHeader(code)
			return
		}

		user := model.NewUser()
		user.IdUser = idUser
		user.Name = name
		user.Phone = phone
		user.Email = email
		user.Password = pass
		user.Foto = photo
		user.Uid = a.UID

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

		resp,err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		if string(resp) == "null"{
			code = http.StatusNotFound
			w.WriteHeader(code)
			return
		}
		code = http.StatusOK
		w.WriteHeader(code)
		w.Write([]byte(resp))
	}
}

func phoneFormat(phone string)string{
	prefix := phone[0]
	if(string(prefix) == "0"){
		val := phone[1:len(phone)]
		phone = fmt.Sprintf("%s%s","+62",val)
	}
	return phone
}

func GetAuthFirebase(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idUser,_ := jsonparser.GetString(body,"IdUser")

		trxRepoImpl := userRepo.NewUserRepositoryImpl(db)
		datas,err := trxRepoImpl.GetAuthFirebase(idUser)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		resp,err := json.Marshal(datas)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		if string(resp) == "null"{
			code = http.StatusNotFound
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
		w.Write(resp)
	}
}

func FindById(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusNotFound
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		idUser,_ := jsonparser.GetString(body,"IdUser") //089123123123

		userRepoImpl := userRepo.NewUserRepositoryImpl(db)

		user,err := userRepoImpl.FindById(idUser)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		resp,_ := json.Marshal(user)
		fmt.Println(string(resp))
		code = http.StatusOK
		w.WriteHeader(code)
		w.Write(resp)
	}
}