package service

import (
	"net/http"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"database/sql"
	genericRepo "github.com/fadhilfcr/oren-service/src/modules/base/repository"
	stringUtil "github.com/fadhilfcr/oren-service/src/util/string"
	constant "github.com/fadhilfcr/oren-service/src/util"
	repo "github.com/fadhilfcr/oren-service/src/modules/rating/repository"
	trxRepo "github.com/fadhilfcr/oren-service/src/modules/transaction/repository"
	"github.com/fadhilfcr/oren-service/src/modules/rating/model"
	"encoding/json"
	"fmt"
)

func SetRating(db *sql.DB)func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		RespChat,_ := jsonparser.GetFloat(body,"RespChat")
		Flexible,_ := jsonparser.GetFloat(body,"Flexible") //089123123123
		Perawatan,_ := jsonparser.GetFloat(body,"Perawatan")
		TepatWaktu,_ := jsonparser.GetFloat(body,"TepatWaktu") //fadhil123
		Friendly,_ := jsonparser.GetFloat(body,"Friendly")

		Desc,_ := jsonparser.GetString(body,"Desc")
		IdTransaksi,_ := jsonparser.GetString(body,"IdTransaksi")
		Flag,_ := jsonparser.GetInt(body,"Flag")

		repoImpl := repo.NewRatingRepositoryImpl(db)
		genericRepoImpl := genericRepo.NewGenericRepositoryImpl(db)

		maxId,_ := genericRepoImpl.FindMaxId(constant.TAG_RATING_PK,constant.TAG_RATING_TBL,constant.TAG_RATING)

		id,err := stringUtil.TableIdFormatter(constant.TAG_RATING,maxId)

		total := RespChat + Flexible + Perawatan +TepatWaktu + Friendly
		idUser,flags,err := repoImpl.FindByIdTrx(IdTransaksi,Flag)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		rating := model.NewRating()
		rating.IdRating = id
		rating.IdTransaksi = IdTransaksi
		rating.IdUser = idUser
		rating.RespChat = RespChat
		rating.Flexible = Flexible
		rating.Perawatan = Perawatan
		rating.TepatWaktu = TepatWaktu
		rating.Friendly = Friendly
		rating.Desc = Desc
		rating.Flag = flags
		rating.Total = total
		err = repoImpl.Save(rating)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		obj,err := repoImpl.FindSummaryByIdUser(idUser)
		if err != nil {
			if err != sql.ErrNoRows {
				w.WriteHeader(code)
				return
			}
		}
		ratingSum := model.NewRatingSummary()
		if err == sql.ErrNoRows {
			ratingSum.IdUser = idUser
			ratingSum.RespChat = RespChat
			ratingSum.Flexible = Flexible
			ratingSum.TepatWaktu = TepatWaktu
			ratingSum.Perawatan = Perawatan
			ratingSum.Friendly = Friendly
			ratingSum.JmlData = 1
			ratingSum.Total = total
			ratingSum.RespChatRat = RespChat
			ratingSum.FlexibleRat = Flexible
			ratingSum.TepatWaktuRat = TepatWaktu
			ratingSum.PerawatanRat= Perawatan
			ratingSum.FriendlyRat = Friendly
			totalRat := total / 5 / 1
			ratingSum.TotalRat = totalRat
		}else{
			respChat := obj.RespChat+RespChat
			flexible := obj.Flexible+Flexible
			perawatan := obj.Perawatan+Perawatan
			tepatWaktu := obj.TepatWaktu+TepatWaktu
			friendly := obj.Friendly+Friendly
			totalRating := respChat + flexible + perawatan +tepatWaktu + friendly
			ratingSum.IdUser = idUser
			ratingSum.RespChat = respChat
			ratingSum.Flexible = flexible
			ratingSum.TepatWaktu = tepatWaktu
			ratingSum.Perawatan = perawatan
			ratingSum.Friendly = friendly
			jmlData := obj.JmlData+1
			ratingSum.JmlData = jmlData
			ratingSum.Total = totalRating
			ratingSum.RespChatRat = float64(respChat)/float64(jmlData)
			ratingSum.FlexibleRat = float64(flexible)/float64(jmlData)
			ratingSum.TepatWaktuRat = float64(tepatWaktu)/float64(jmlData)
			ratingSum.PerawatanRat = float64(perawatan)/float64(jmlData)
			ratingSum.FriendlyRat = float64(friendly)/float64(jmlData)
			ratingSum.TotalRat = float64(totalRating)/5/float64(jmlData)
		}
		err = repoImpl.SaveSum(ratingSum)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		trxRepo := trxRepo.NewTrxRepositoryImpl(db);
		trxObj,err := trxRepo.FindById(IdTransaksi)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		var idStatus string
		if(trxObj.IdStatus == constant.STS_TRX_RATING_BOTH_10){
			if(Flag == 0){
				idStatus = constant.STS_TRX_RATING_RENT_11
			}else{
				idStatus = constant.STS_TRX_RATING_OWNER_12
			}
		}else if(trxObj.IdStatus == constant.STS_TRX_RATING_RENT_11 || trxObj.IdStatus == constant.STS_TRX_RATING_OWNER_12){
			idStatus = constant.STS_TRX_FINISHED_BOTH_13
		}else{
			idStatus = constant.STS_TRX_RATING_BOTH_10
		}
		err = trxRepo.UpdateStatus(IdTransaksi,idStatus)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
	}
}

func FindAccountByIdUser(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		IdUser,_ := jsonparser.GetString(body,"IdUser")

		trxRepoImpl := repo.NewRatingRepositoryImpl(db)
		datas,err := trxRepoImpl.FindAccountByIdUser(IdUser)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(code)
			return
		}
		resp,err := json.Marshal(datas)
		if err != nil {
			fmt.Println(err.Error())

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

func FindComments(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		IdUser,_ := jsonparser.GetString(body,"IdUser")

		trxRepoImpl := repo.NewRatingRepositoryImpl(db)
		datas,err := trxRepoImpl.FindComments(IdUser)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(code)
			return
		}
		resp,err := json.Marshal(datas)
		if err != nil {
			fmt.Println(err.Error())
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

