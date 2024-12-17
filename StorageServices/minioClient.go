package StorageServices

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"mime/multipart"
	"time"
)

type minioClient struct {
	client *minio.Client
}

func NewMinioClient(host string, port string, accessKey string, secretKey string) *minioClient {
	return &minioClient{client: initiateClient(host, port, accessKey, secretKey)}
}

func initiateClient(host string, port string, accessKey string, secretKey string) *minio.Client {
	
	mc, err := minio.New(host+":"+port, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	
	if err != nil {
		log.Fatalln(err)
	}
	
	return mc
	
}

func (c *minioClient) UploadObject(bucketName, objectName string, file multipart.File, header *multipart.FileHeader) error {
	
	contentType := header.Header.Get("Content-Type")
	
	_, err := c.client.PutObject(context.TODO(), bucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	
	if err != nil {
		log.Printf("Failed to upload %s: %v", objectName, err)
		return err
	}
	return nil
}

func (c *minioClient) DownloadObject(bucketName, objectName string) ([]byte, error) {
	object, err := c.client.GetObject(context.TODO(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer func(object *minio.Object) {
		errClose := object.Close()
		if errClose != nil {
			log.Println(errClose)
		}
	}(object)
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, object); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *minioClient) DownloadObjectStream(bucketName, objectName string) (io.Reader, error) {
	object, err := c.client.GetObject(context.TODO(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	
	return object, nil
}

// in a real production code i will not implement this like this
func (c *minioClient) GeneratePresignedURL(method string, bucketName string, objectName string) (string, error) {
	// the expiry time for the presigned URL
	expiry := 5 * time.Minute
	
	if method == "PUT" {
		presignedURL, err := c.client.PresignedPutObject(context.TODO(), bucketName, objectName, expiry)
		if err != nil {
			log.Printf("Failed to generate PUT presigned URL for %s in bucket %s: %v", objectName, bucketName, err)
			return "", err
		}
		return presignedURL.String(), nil
		
	} else if method == "GET" {
		presignedURL, err := c.client.PresignedGetObject(context.TODO(), bucketName, objectName, expiry, nil)
		if err != nil {
			log.Printf("Failed to generate GET presigned URL for %s in bucket %s: %v", objectName, bucketName, err)
			return "", err
		}
		return presignedURL.String(), nil
	}
	return "", nil
}
