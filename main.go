package main

import (
	"test/StorageServices"
	http "test/engine"
)

func main() {

	minioClient := StorageServices.NewMinioClient("localhost", "9000", "J3KmVTuNp6w0LtpmXQuN", "pxqckfe91DqpvAJcrr6Fwqut1XEPeBGOgHxLObV0")
	// I Can also use AWSClient
	//awsClient := StorageServices.NewAWSS3Client("http://localhost:9000", "us-east-1", "J3KmVTuNp6w0LtpmXQuN", "pxqckfe91DqpvAJcrr6Fwqut1XEPeBGOgHxLObV0", true)
	router := http.NewRouter(minioClient)

	ginEngine := http.NewGinEngine("8080", router.InstallFileHandler)
	ginEngine.RunHttpServer()

}
