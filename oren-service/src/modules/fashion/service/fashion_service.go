package service

import (
	"database/sql"
	"net/http"
	"io/ioutil"
	"fmt"
	"time"
	"os"
	"io"
	"mime/multipart"
	constant "github.com/fadhilfcr/oren-service/src/util"
	stringUtil "github.com/fadhilfcr/oren-service/src/util/string"
	genericRepo "github.com/fadhilfcr/oren-service/src/modules/base/repository"
	repo "github.com/fadhilfcr/oren-service/src/modules/fashion/repository"
	model "github.com/fadhilfcr/oren-service/src/modules/fashion/model"
	shippmentService "github.com/fadhilfcr/oren-service/src/modules/shippment/service"
	"strconv"
	"github.com/minio/minio-go"
	"net/url"
	"strings"
	"github.com/buger/jsonparser"
	"encoding/json"
)

func FindDataById(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("find data by id")
		code := http.StatusNotFound
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		id,_ := jsonparser.GetString(body,"IdBusana")
		if len(id) <= 0 {
			w.WriteHeader(code)
			return
		}

		repoImpl := repo.NewFashionRepositoryImpl(db)

		var objs *model.FashionDetail
		objs, err = repoImpl.FindById(id)
		if err != nil {
			fmt.Println("get Data : "+err.Error())
			w.WriteHeader(code)
			return
		}

		objs.Pengiriman,err = shippmentService.FindByIdBusana(db,id)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		resp,err := json.Marshal(objs)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		if string(resp) == "null"{
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
		w.Write(resp)
	}
}


func FindAllFashionByParameter(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Println("find data by param")
		code := http.StatusNotFound
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		idUser,_ := jsonparser.GetString(body,"IdUser")
		name,_ := jsonparser.GetString(body,"TxtSearch")
		idCategory,_ := jsonparser.GetString(body,"IdCategory")
		location,_ := jsonparser.GetString(body,"Location")
		idShippment,_ := jsonparser.GetString(body,"IdShippment")
		startDate,_ := jsonparser.GetInt(body,"StartDate")
		endDate,_ := jsonparser.GetInt(body,"EndDate")

		param := make(map[string]string)
		var objs *model.FashionSearchs

		if(len(idUser) > 0){
			param["idUser"] = strings.Trim(idUser," ")
		}

		if(len(name) > 0){
			param["name"] = strings.Trim(name," ")
		}

		if(len(idCategory) > 0){
			param["idCategory"] = strings.Trim(idCategory," ")
		}

		if(len(location) > 0){
			param["location"] = strings.Trim(location," ")
		}

		if(len(idShippment) > 0){
			param["idShippment"] = strings.Trim(idShippment," ")
		}

		if(startDate > 0){
			param["startDate"] = strconv.FormatInt(startDate,10)
		}

		if(endDate > 0){
			param["endDate"] = strconv.FormatInt(endDate,10)
		}

		repoImpl := repo.NewFashionRepositoryImpl(db)

		objs, err = repoImpl.FindAllByParameter(param)
		if err != nil {
			fmt.Println("get Data : "+err.Error())
			w.WriteHeader(code)
			return
		}

		resp,err := json.Marshal(objs)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		if string(resp) == "null"{
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
		w.Write(resp)
	}
}

func AddFashion(db *sql.DB,minioClient *minio.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("service.fashion.addFashion")
		code := http.StatusInternalServerError
		message := `{"message":"Service Not Available"}`
		fileMain,headerMain,err := r.FormFile("fileMain")
		file1,header1,err := r.FormFile("file1")
		file2,header2,err := r.FormFile("file2")
		file3,header3,err := r.FormFile("file3")
		file4,header4,err := r.FormFile("file4")

		if err != nil {
			fmt.Println("formFile : "+err.Error())
			w.WriteHeader(code)
			return
		}

		idUser := r.FormValue("IdUser")
		name := r.FormValue("name")
		idCategory := r.FormValue("idCategory")
		weight := r.FormValue("weight")
		deposit := r.FormValue("deposit")
		price := r.FormValue("price")
		shippment := r.FormValue("shippment")
		desc := r.FormValue("desc")

		if len(headerMain.Filename) <= 0 || len(header1.Filename) <= 0 || len(header2.Filename) <= 0 || len(header3.Filename) <= 0 || len(header4.Filename) <= 0 {
			fmt.Println("formFile : "+err.Error())
			code = http.StatusBadRequest
			message = fmt.Sprintf(`{"message":"Error : %s"}`,"Payload is blank or not complete")
			w.WriteHeader(code)
			w.Write([]byte(message))
			return
		}

		genericRepoImpl := genericRepo.NewGenericRepositoryImpl(db)
		repoImpl := repo.NewFashionRepositoryImpl(db)

		maxId,_ := genericRepoImpl.FindMaxId(constant.TAG_BUSANA_PK,constant.TAG_BUSANA_TBL,constant.TAG_BUSANA)

		idBusana,_ := stringUtil.TableIdFormatter(constant.TAG_BUSANA,maxId)

		flWeight,_ := strconv.ParseFloat(weight,10)
		flPrice,_ := strconv.ParseFloat(price,10)
		flDeposit,_ := strconv.ParseFloat(deposit,10)

		fashion := model.NewFashion()
		fashion.IdBusana = idBusana
		fashion.NamaBusana = name
		fashion.IdKategori = idCategory
		fashion.IdUser = idUser
		fashion.Berat = flWeight
		fashion.Harga = flPrice
		fashion.Deposit = flDeposit
		fashion.Deskripsi = desc

		urlMain,filenameMain,err := saveMinio(idBusana,headerMain,fileMain,minioClient)
		if err != nil {
			fmt.Println("save Minio : "+err.Error())
			w.WriteHeader(code)
			return
		}
		url1,filename1,err := saveMinio(idBusana,header1,file1,minioClient)
		if err != nil {
			fmt.Println("save Minio : "+err.Error())
			w.WriteHeader(code)
			return
		}
		url2,filename2,err := saveMinio(idBusana,header2,file2,minioClient)
		if err != nil {
			fmt.Println("save Minio : "+err.Error())
			w.WriteHeader(code)
			return
		}
		url3,filename3,err := saveMinio(idBusana,header3,file3,minioClient)
		if err != nil {
			fmt.Println("save Minio : "+err.Error())
			w.WriteHeader(code)
			return
		}
		url4,filename4,err := saveMinio(idBusana,header4,file4,minioClient)
		if err != nil {
			fmt.Println("save Minio : "+err.Error())
			w.WriteHeader(code)
			return
		}

		fashion.FotoUtama = filenameMain
		fashion.FotoSatu = filename1
		fashion.FotoDua = filename2
		fashion.FotoTiga = filename3
		fashion.FotoEmpat = filename4

		fashion.UrlFotoUtama = urlMain
		fashion.UrlFotoSatu = url1
		fashion.UrlFotoDua = url2
		fashion.UrlFotoTiga = url3
		fashion.UrlFotoEmpat = url4
		fashion.IdStatus = constant.STS_FASHION_NEW_01

		err = repoImpl.Save(fashion)
		if err != nil {
			fmt.Println("save Data : "+err.Error())
			w.WriteHeader(code)
			return
		}

		err = shippmentService.SaveShippmentFashion(db,shippment,idBusana)
		if err != nil {
			fmt.Println("save Data Shippment Fashion: "+err.Error())
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
		w.Write([]byte(message))
	}
}

func createTempFile(file multipart.File,header *multipart.FileHeader,unixTime int64)(error){
	filename := fmt.Sprintf("%d_%s", unixTime, header.Filename)
	folder := "temp/"
	tmpPath := fmt.Sprintf("%s%s", folder, filename)
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return err
}

func saveMinio(idBusana string ,header *multipart.FileHeader, file multipart.File , minioClient *minio.Client)(string,string,error){
	var url *url.URL
	contentType := header.Header.Get("Content-Type")
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	folder := "temp/"
	tmpPath := fmt.Sprintf("%s%s", folder, filename)
	f, err := os.Create(tmpPath)
	if err != nil {
		return "","",err
	}
	defer f.Close()
	io.Copy(f, file)

	bucketName := strings.ToLower(idBusana)
	existBucket,err := minioClient.BucketExists(bucketName)
	if(!existBucket){
		minioClient.MakeBucket(bucketName,"ap-southeast-1")
	}

	_, err = minioClient.FPutObject(bucketName, filename, tmpPath, minio.PutObjectOptions{ContentType: contentType})
	if err == nil {
		url, err = minioClient.PresignedGetObject(bucketName, filename, time.Second*24*60*60*7, nil)
		if err != nil {
			return "","",err
		}
	}
	os.Remove(tmpPath)
	return url.String(),filename,nil
}

func FindByIdUser(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Println("find data by param")
		code := http.StatusNotFound
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		idUser,_ := jsonparser.GetString(body,"IdUser")


		repoImpl := repo.NewFashionRepositoryImpl(db)

		objs, err := repoImpl.FindByIdUser(idUser)
		if err != nil {
			fmt.Println("get Data : "+err.Error())
			w.WriteHeader(code)
			return
		}

		resp,err := json.Marshal(objs)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		if string(resp) == "null"{
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
		w.Write(resp)
	}
}