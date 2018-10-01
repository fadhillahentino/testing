package main

import(
	"github.com/fadhilfcr/oren-service/config"
	"fmt"
	"database/sql"
	"net/http"
	userService "github.com/fadhilfcr/oren-service/src/modules/users/service"
	fashionService "github.com/fadhilfcr/oren-service/src/modules/fashion/service"
	genderService "github.com/fadhilfcr/oren-service/src/modules/gender/service"
	categoryService "github.com/fadhilfcr/oren-service/src/modules/category/service"
	shippmentService "github.com/fadhilfcr/oren-service/src/modules/shippment/service"
	locationService "github.com/fadhilfcr/oren-service/src/modules/province/service"
	addressService "github.com/fadhilfcr/oren-service/src/modules/address/service"
	trxService "github.com/fadhilfcr/oren-service/src/modules/transaction/service"
	accService "github.com/fadhilfcr/oren-service/src/modules/account/service"
	ratService "github.com/fadhilfcr/oren-service/src/modules/rating/service"
	"github.com/minio/minio-go"
	"golang.org/x/net/context"
)

const PORT_NUMBER = 8082

func conn() (*sql.DB){
	db,err := config.GetMysqlDB()
	if err != nil{
		fmt.Println(err.Error())
	}
	return db
}

func connMinio() (*minio.Client){
	minio := config.GetMinio()
	return minio
}

func phoneFormat(phone string)string{
	prefix := phone[0]
	if(string(prefix) == "0"){
		val := phone[1:len(phone)]
		phone = fmt.Sprintf("%s%s","+62",val)
	}
	return phone
}

func main() {
	fmt.Println("==== Main services ====")
	fmt.Println(fmt.Sprintf(" Run on Port :%d", PORT_NUMBER))
	fmt.Println("======================")
	db := conn()
	defer db.Close()
	minio := connMinio()
	mux := http.NewServeMux()
	mux.HandleFunc("/oren/user/login", login(db))
	mux.HandleFunc("/oren/user/registration", registration(db))
	mux.HandleFunc("/oren/user/findById", FindById(db))
	mux.HandleFunc("/oren/user/getAuthFirebase", GetAuthFirebase(db))
	mux.HandleFunc("/oren/fashion/findDataById", findDataById(db))
	mux.HandleFunc("/oren/fashion/findFashionByParameter", findFashionByParameter(db))
	mux.HandleFunc("/oren/fashion/findByIdUser", findByIdUser(db))
	mux.HandleFunc("/oren/fashion/addFashion", addFashion(db,minio))
	mux.HandleFunc("/oren/gender/findAll", findAllGender(db))
	mux.HandleFunc("/oren/category/findAll", findAllCategory(db))
	mux.HandleFunc("/oren/category/findByGender", findCategoryByGender(db))
	mux.HandleFunc("/oren/shippment/findAll", findAllShippment(db))
	mux.HandleFunc("/oren/shippment/findByIdBusana", FindAllByIdBusana(db))
	mux.HandleFunc("/oren/province/findAll", findAllProvince(db))
	mux.HandleFunc("/oren/city/findAll", findAllCity(db))
	mux.HandleFunc("/oren/address/findByIdUser", findAddressByIdUser(db))
	mux.HandleFunc("/oren/trx/updateStatus", UpdateStatus(db))
	mux.HandleFunc("/oren/trx/trxConfirmRent", TrxConfirmRent(db))
	mux.HandleFunc("/oren/trx/ListNotificationUser", ListNotificationUser(db))
	mux.HandleFunc("/oren/trx/rentDetail", RentDetail(db))
	mux.HandleFunc("/oren/trx/approveTrx", ApproveTrx(db))
	mux.HandleFunc("/oren/trx/confirmBilling", ConfirmBilling(db))
	mux.HandleFunc("/oren/trx/confirmBillingSave", ConfirmBillingSave(db,minio))
	mux.HandleFunc("/oren/trx/receiptPayment", ReceiptPayment(db))
	mux.HandleFunc("/oren/trx/sendProduct", SendProduct(db))
	mux.HandleFunc("/oren/trx/sendProductSave", SendProductSave(db,minio))
	mux.HandleFunc("/oren/trx/receiveProduct", ReceiveProduct(db))
	mux.HandleFunc("/oren/trx/retriveProduct", RetriveProduct(db))
	mux.HandleFunc("/oren/trx/retriveProductSave", RetriveProductSave(db,minio))
	mux.HandleFunc("/oren/trx/receiveProductOwner", ReceiveProductOwner(db))
	mux.HandleFunc("/oren/rating/save", SetRating(db))
	mux.HandleFunc("/oren/rating/findAccountByIdUser", FindAccountByIdUser(db))
	mux.HandleFunc("/oren/rating/findComments", FindComments(db))
	mux.HandleFunc("/oren/account/findAccountById", FindAccountById(db))
	http.ListenAndServe(fmt.Sprintf(":%d", PORT_NUMBER), mux)
}


//main
func login(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	return userService.Login_service(db)
}

func registration(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	return userService.Registration_service(db,context.Background())
}

func FindById(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	return userService.FindById(db)
}

func GetAuthFirebase(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.GetAuthFirebase")
	return userService.GetAuthFirebase(db)
}

//fashion
func findDataById(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	return fashionService.FindDataById(db)
}

func findFashionByParameter(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	return fashionService.FindAllFashionByParameter(db)
}

func findByIdUser(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	return fashionService.FindByIdUser(db)
}

func addFashion(db *sql.DB,minioClient *minio.Client) func(w http.ResponseWriter, r *http.Request){
	return fashionService.AddFashion(db,minioClient)
}

//gender
func findAllGender(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	return genderService.FindAll(db)
}

//category
func findAllCategory(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	return categoryService.FindAll(db)
}

func findCategoryByGender(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	return categoryService.FindbyGender(db)
}

//shippment
func findAllShippment(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.shippment.findAll")
	return shippmentService.FindAll(db)
}

func FindAllByIdBusana(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.shippment.FindAllByIdBusana")
	return shippmentService.FindAllByIdBusana(db)
}

//Location
func findAllProvince(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.location.findAllProvince")
	return locationService.FindAllProvince(db)
}

func findAllCity(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.location.findAllCity")
	return locationService.FindAllCity(db)
}

func findAddressByIdUser(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.location.findAllCity")
	return addressService.FindAddressByIdUser(db)
}

//Transaction
func UpdateStatus(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.UpdateStatus")
	return trxService.UpdateStatus(db)
}

func TrxConfirmRent(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.Trx_confirm_rent")
	return trxService.Trx_confirm_rent(db)
}

func ListNotificationUser(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.ListNotificationUser")
	return trxService.ListNotificationUser(db)
}
func RentDetail(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.RentDetail")
	return trxService.RentDetail(db)
}

func ApproveTrx(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.ApproveTrx")
	return trxService.ApproveTrx(db)
}

func ConfirmBilling(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.ConfirmBilling")
	return trxService.ConfirmBilling(db)
}

func ConfirmBillingSave(db *sql.DB,minioClient *minio.Client) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.ConfirmBillingSave")
	return trxService.ConfirmBillingSave(db,minioClient)
}

func ReceiptPayment(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.ReceiptPayment")
	return trxService.ReceiptPayment(db)
}

func SendProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.SendProduct")
	return trxService.SendProduct(db)
}

func SendProductSave(db *sql.DB,minioClient *minio.Client) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.SendProductSave")
	return trxService.SendProductSave(db,minioClient)
}

func ReceiveProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.ReceiveProduct")
	return trxService.ReceiveProduct(db)
}

func RetriveProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.RetriveProduct")
	return trxService.RetriveProduct(db)
}

func RetriveProductSave(db *sql.DB,minioClient *minio.Client) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.RetriveProductSave")
	return trxService.RetriveProductSave(db,minioClient)
}

func ReceiveProductOwner(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.ReceiveProductOwner")
	return trxService.ReceiveProductOwner(db)
}

func SetRating(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.ReceiveProductOwner")
	return ratService.SetRating(db)
}

func FindAccountByIdUser(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.findAccountByIdUser")
	return ratService.FindAccountByIdUser(db)
}

func FindComments(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.FindComments")
	return ratService.FindComments(db)
}

//Account
func FindAccountById(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
	fmt.Println("main.trx.FindAccountById")
	return accService.FindAccountById(db)
}
