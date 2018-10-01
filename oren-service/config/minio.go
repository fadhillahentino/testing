package config

import (
	"github.com/minio/minio-go"
	"log"
)

func GetMinio() *minio.Client{
	endpoint:="128.199.129.191:9000"
	accessKeyID := "5TKTJZDVIWJVZKDMLAS1"
	secretAccessKey := "Y6BU0zP0pTLQEvUF5SeIjxT+NwGLCCNJySc3G65w"
	useSSL := false
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
		return nil
	}else{
		return minioClient
	}
}
