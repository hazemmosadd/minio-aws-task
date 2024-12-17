package StorageServices

import (
	"io"
	"mime/multipart"
)

type StorageServiceClient interface {
	UploadObject(bucketName, objectName string, file multipart.File, header *multipart.FileHeader) error
	DownloadObject(bucketName string, objectKey string) ([]byte, error)
	DownloadObjectStream(bucketName string, objectKey string) (io.Reader, error)
	GeneratePresignedURL(method string, bucketName string, objectName string) (string, error)
}
