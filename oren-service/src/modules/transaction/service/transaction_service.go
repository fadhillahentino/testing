package service

import (
	"net/http"
	"io/ioutil"
	"github.com/buger/jsonparser"
	"database/sql"
	genericRepo "github.com/fadhilfcr/oren-service/src/modules/base/repository"

	stringUtil "github.com/fadhilfcr/oren-service/src/util/string"
	constant "github.com/fadhilfcr/oren-service/src/util"
	repo "github.com/fadhilfcr/oren-service/src/modules/transaction/repository"
	"github.com/fadhilfcr/oren-service/src/modules/transaction/model"

	addrModel "github.com/fadhilfcr/oren-service/src/modules/address/model"
	addrRepo "github.com/fadhilfcr/oren-service/src/modules/address/repository"
	accModel "github.com/fadhilfcr/oren-service/src/modules/account/model"
	accRepo "github.com/fadhilfcr/oren-service/src/modules/account/repository"
	fashionRepo "github.com/fadhilfcr/oren-service/src/modules/fashion/repository"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go"
	"mime/multipart"
	"net/url"
	"time"
	"os"
	"io"
	"strings"
	"strconv"
)

func UpdateStatus(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idTransaksi,_ := jsonparser.GetString(body,"IdTransaksi")
		idStatus,_ := jsonparser.GetString(body,"IdStatus")

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		err = trxRepoImpl.UpdateStatus(idTransaksi,idStatus)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
	}
}

func Trx_confirm_rent(db *sql.DB)func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idBusana,_ := jsonparser.GetString(body,"IdBusana")
		idPemilik,_ := jsonparser.GetString(body,"IdPemilik")
		idPenyewa,_ := jsonparser.GetString(body,"IdPenyewa") //089123123123
		startDate,_ := jsonparser.GetInt(body,"StartDate")
		endDate,_ := jsonparser.GetInt(body,"EndDate")
		totalHarga,_ := jsonparser.GetFloat(body,"TotalHarga")
		idPengiriman,_ := jsonparser.GetString(body,"IdPengiriman")
		idStatus,_ := jsonparser.GetString(body,"IdStatus")

		idAlamat,_ := jsonparser.GetString(body,"IdAlamat")
		idProvinsi,_ := jsonparser.GetString(body,"IdProvinsi") //fadhil123
		idKabupaten,_ := jsonparser.GetString(body,"IdKabupaten")
		kecamatan,_ := jsonparser.GetString(body,"Kecamatan")
		kodePos,_ := jsonparser.GetString(body,"KodePos")
		alamat,_ := jsonparser.GetString(body,"Alamat")

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		addrRepoImpl := addrRepo.NewAddressRepositoryImpl(db)
		genericRepoImpl := genericRepo.NewGenericRepositoryImpl(db)
		fashionRepoImpl := fashionRepo.NewFashionRepositoryImpl(db)

		maxId,_ := genericRepoImpl.FindMaxId(constant.TAG_TRANSAKSI_PK,constant.TAG_TRANSAKSI_TBL,constant.TAG_TRANSAKSI)
		idTransaksi,err := stringUtil.TableIdFormatter(constant.TAG_TRANSAKSI,maxId)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		obj := model.NewTransaction()
		obj.IdBusana = idBusana
		obj.IdTransaksi = idTransaksi
		obj.IdPemilik = idPemilik
		obj.IdPenyewa = idPenyewa
		obj.StartDate = startDate
		obj.EndDate = endDate
		obj.TotalHarga = totalHarga
		obj.IdPengiriman = idPengiriman
		obj.IdStatus = idStatus

		addr := addrModel.NewAddress()
		addr.IdUser = idPenyewa
		addr.Alamat = alamat
		addr.Provinsi = idProvinsi
		addr.Kabupaten = idKabupaten
		addr.Kecamatan = kecamatan
		addr.KodePos = kodePos

		if(len(idAlamat) <= 0) {
			maxIdAlamat, _ := genericRepoImpl.FindMaxId(constant.TAG_ADDRESS_PK, constant.TAG_ADDRESS_TBL, constant.TAG_ADDRESS)
			idAlamat,err = stringUtil.TableIdFormatter(constant.TAG_ADDRESS, maxIdAlamat)
			if err != nil {
				w.WriteHeader(code)
				return
			}
			addr.IdAlamat = idAlamat
			err = addrRepoImpl.Save(addr)
			if err != nil {
				w.WriteHeader(code)
				return
			}
		}else{
			err = addrRepoImpl.Update(idAlamat,addr)
			if err != nil {
				w.WriteHeader(code)
				return
			}
		}
		obj.IdAlamat = idAlamat
		err = trxRepoImpl.Save(obj)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		err = trxRepoImpl.SaveDetail(obj)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		err = fashionRepoImpl.UpdateStatus(idBusana,constant.STS_FASHION_RENT_02)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
	}
}

func ListNotificationUser(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idUser,_ := jsonparser.GetString(body,"IdUser")

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		datas,err := trxRepoImpl.ListNotificationUser(idUser)
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

func RentDetail(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idTransaksi,_ := jsonparser.GetString(body,"IdTransaksi")

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		datas,err := trxRepoImpl.RentDetail(idTransaksi)
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

func ApproveTrx(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idTransaksi,_ := jsonparser.GetString(body,"IdTransaksi")
		idStatus,_ := jsonparser.GetString(body,"IdStatus")
		idUser,_ := jsonparser.GetString(body,"IdUser")
		idRekening,_ := jsonparser.GetString(body,"IdRekening")
		nomorRekening,_ := jsonparser.GetString(body,"NomorRekening")
		namaBank,_ := jsonparser.GetString(body,"NamaBank")
		atasNama,_ := jsonparser.GetString(body,"AtasNama")
		biayaPengiriman,_ := jsonparser.GetFloat(body,"BiayaPengiriman")

		idAlamat,_ := jsonparser.GetString(body,"IdAlamat")
		idProvinsi,_ := jsonparser.GetString(body,"IdProvinsi") //fadhil123
		idKabupaten,_ := jsonparser.GetString(body,"IdKabupaten")
		kecamatan,_ := jsonparser.GetString(body,"Kecamatan")
		kodePos,_ := jsonparser.GetString(body,"KodePos")
		alamat,_ := jsonparser.GetString(body,"Alamat")

		addrRepoImpl := addrRepo.NewAddressRepositoryImpl(db)
		accRepoImpl := accRepo.NewAccountRepositoryImpl(db)
		genericRepoImpl := genericRepo.NewGenericRepositoryImpl(db)

		acc := accModel.NewAccount()
		acc.AtasNama = atasNama
		acc.NamaBank = namaBank
		acc.NomorRekening = nomorRekening
		acc.IdUser = idUser


		if(len(idRekening) <= 0){
			maxId,_ := genericRepoImpl.FindMaxId(constant.TAG_REKENING_PK,constant.TAG_REKENING_TBL,constant.TAG_REKENING)
			idRekening,err = stringUtil.TableIdFormatter(constant.TAG_REKENING,maxId)
			if err != nil {
				w.WriteHeader(code)
				return
			}
			acc.IdRekening = idRekening
			err = accRepoImpl.Save(acc)
			if err != nil {
				w.WriteHeader(code)
				return
			}
		}else{
			err = accRepoImpl.Update(idRekening,acc)
			if err != nil {
				w.WriteHeader(code)
				return
			}
		}

		addr := addrModel.NewAddress()
		addr.IdUser = idUser
		addr.Alamat = alamat
		addr.Provinsi = idProvinsi
		addr.Kabupaten = idKabupaten
		addr.Kecamatan = kecamatan
		addr.KodePos = kodePos

		if(len(idAlamat) <= 0) {
			maxIdAlamat, _ := genericRepoImpl.FindMaxId(constant.TAG_ADDRESS_PK, constant.TAG_ADDRESS_TBL, constant.TAG_ADDRESS)
			idAlamat,err = stringUtil.TableIdFormatter(constant.TAG_ADDRESS, maxIdAlamat)
			if err != nil {
				w.WriteHeader(code)
				return
			}
			addr.IdAlamat = idAlamat
			err = addrRepoImpl.Save(addr)
			if err != nil {
				w.WriteHeader(code)
				return
			}
		}else{
			err = addrRepoImpl.Update(idAlamat,addr)
			if err != nil {
				w.WriteHeader(code)
				return
			}
		}
		idAlamatKembali := idAlamat

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		err = trxRepoImpl.UpdateStatusApprove(idTransaksi,idStatus,idRekening,biayaPengiriman,idAlamatKembali)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		code = http.StatusOK
		w.WriteHeader(code)
	}
}

func ConfirmBilling(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idTransaksi,_ := jsonparser.GetString(body,"IdTransaksi")

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		datas,err := trxRepoImpl.ConfirmBilling(idTransaksi)
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

func ConfirmBillingSave(db *sql.DB,minioClient *minio.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("service.fashion.addFashion")
		code := http.StatusInternalServerError
		message := `{"message":"Service Not Available"}`
		fileMain, headerMain, err := r.FormFile("fileMain")

		if err != nil {
			fmt.Println("formFile : " + err.Error())
			w.WriteHeader(code)
			return
		}

		idTrx := r.FormValue("idTransaksi")
		accName := r.FormValue("accName")
		bank := r.FormValue("bank")
		nominal := r.FormValue("nominal")

		if len(headerMain.Filename) <= 0{
			fmt.Println("formFile : " + err.Error())
			code = http.StatusBadRequest
			message = fmt.Sprintf(`{"message":"Error : %s"}`, "Payload is blank or not complete")
			w.WriteHeader(code)
			w.Write([]byte(message))
			return
		}

		urlMain,filenameMain,err := saveMinio(idTrx,headerMain,fileMain,minioClient)
		if err != nil {
			fmt.Println("save Minio : "+err.Error())
			w.WriteHeader(code)
			return
		}

		trxRent := model.NewTransaction()
		trxRent.NamaBank = bank
		trxRent.AtasNama = accName
		trxRent.NominaTransfer,_ = strconv.ParseFloat(nominal,10)
		trxRent.BuktiTransfer = filenameMain
		trxRent.UrlBuktiTransfer = urlMain
		trxRent.IdTransaksi = idTrx
		trxRent.IdStatus = constant.STS_TRX_CONFIRMATIONPAYMENT_OWNER_04
		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		err = trxRepoImpl.ConfirmBillingSave(trxRent)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		code = http.StatusOK
		w.WriteHeader(code)
	}
}

func saveMinio(idTrx string ,header *multipart.FileHeader, file multipart.File , minioClient *minio.Client)(string,string,error){
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

	bucketName := strings.ToLower(idTrx)
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

func ReceiptPayment(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idTransaksi,_ := jsonparser.GetString(body,"IdTransaksi")

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		datas,err := trxRepoImpl.ReceiptPayment(idTransaksi)
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

func SendProduct(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idTransaksi,_ := jsonparser.GetString(body,"IdTransaksi")

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		datas,err := trxRepoImpl.SendProduct(idTransaksi)
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

func SendProductSave(db *sql.DB,minioClient *minio.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("service.fashion.addFashion")
		code := http.StatusInternalServerError
		message := `{"message":"Service Not Available"}`
		fileMain, headerMain, err := r.FormFile("fileMain")

		if err != nil {
			fmt.Println("formFile : " + err.Error())
			w.WriteHeader(code)
			return
		}

		idTrx := r.FormValue("idTransaksi")
		resi := r.FormValue("resi")

		if len(headerMain.Filename) <= 0{
			fmt.Println("formFile : " + err.Error())
			code = http.StatusBadRequest
			message = fmt.Sprintf(`{"message":"Error : %s"}`, "Payload is blank or not complete")
			w.WriteHeader(code)
			w.Write([]byte(message))
			return
		}

		urlMain,filenameMain,err := saveMinio(idTrx,headerMain,fileMain,minioClient)
		if err != nil {
			fmt.Println("save Minio : "+err.Error())
			w.WriteHeader(code)
			return
		}

		trxRent := model.NewTransaction()
		trxRent.Resi = resi
		trxRent.IdTransaksi = idTrx
		trxRent.BuktiPengiriman = filenameMain
		trxRent.UrlBuktiPengiriman = urlMain
		trxRent.IdStatus = constant.STS_TRX_SHIPPED_OWNER_07
		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		err = trxRepoImpl.SendProductSave(trxRent)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		code = http.StatusOK
		w.WriteHeader(code)
	}
}

func ReceiveProduct(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idTransaksi,_ := jsonparser.GetString(body,"IdTransaksi")

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		datas,err := trxRepoImpl.ReceiveProduct(idTransaksi)
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

func RetriveProduct(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idTransaksi,_ := jsonparser.GetString(body,"IdTransaksi")

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		datas,err := trxRepoImpl.RetriveProduct(idTransaksi)
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

func RetriveProductSave(db *sql.DB,minioClient *minio.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("service.fashion.addFashion")
		code := http.StatusInternalServerError
		message := `{"message":"Service Not Available"}`
		fileMain, headerMain, err := r.FormFile("fileMain")

		if err != nil {
			fmt.Println("formFile : " + err.Error())
			w.WriteHeader(code)
			return
		}

		idTrx := r.FormValue("idTransaksi")
		resi := r.FormValue("resi")
		idPengirimanKembali := r.FormValue("idPengirimanKembali")

		if len(headerMain.Filename) <= 0{
			fmt.Println("formFile : " + err.Error())
			code = http.StatusBadRequest
			message = fmt.Sprintf(`{"message":"Error : %s"}`, "Payload is blank or not complete")
			w.WriteHeader(code)
			w.Write([]byte(message))
			return
		}

		urlMain,_,err := saveMinio(idTrx,headerMain,fileMain,minioClient)
		if err != nil {
			fmt.Println("save Minio : "+err.Error())
			w.WriteHeader(code)
			return
		}

		trxRent := model.NewTransaction()
		trxRent.IdPengirimanKembali = idPengirimanKembali
		trxRent.IdTransaksi = idTrx
		trxRent.BuktiPengirimanKembali = resi
		trxRent.UrlBuktiPengirimanKembali = urlMain
		trxRent.IdStatus = constant.STS_TRX_RETRIVE_RENT_09
		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		err = trxRepoImpl.RetriveProductSave(trxRent)
		if err != nil {
			w.WriteHeader(code)
			return
		}
		code = http.StatusOK
		w.WriteHeader(code)
	}
}

func ReceiveProductOwner(db *sql.DB)func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(code)
			return
		}

		idTransaksi,_ := jsonparser.GetString(body,"IdTransaksi")

		trxRepoImpl := repo.NewTrxRepositoryImpl(db)
		datas,err := trxRepoImpl.ReceiveProductOwner(idTransaksi)
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