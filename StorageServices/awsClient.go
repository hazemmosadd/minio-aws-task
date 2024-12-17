package StorageServices

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"mime/multipart"
	"time"
)

type awsS3Client struct {
	client *s3.S3
}

func NewAWSS3Client(endpoint, region, accessKey, secretKey string, usePathStyle bool) *awsS3Client {
	return &awsS3Client{client: initiateAWSClient(endpoint, region, accessKey, secretKey)}
}

func initiateAWSClient(endpoint, region, accessKey, secretKey string) *s3.S3 {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true), // Required for MinIO
	})
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}

	svc := s3.New(sess)
	return svc

}

func (c *awsS3Client) UploadObject(bucketName, objectName string, file multipart.File, header *multipart.FileHeader) error {

	_, err := c.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectName),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})

	if err != nil {
		log.Printf("Failed to upload %s: %v", objectName, err)
		return err
	}
	return nil
}

func (c *awsS3Client) DownloadObject(bucketName string, objectKey string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}
	result, err := c.client.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, result.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *awsS3Client) DownloadObjectStream(bucketName string, objectKey string) (io.Reader, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}
	result, err := c.client.GetObject(input)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (c *awsS3Client) GeneratePresignedURL(method string, bucketName string, objectName string) (string, error) {

	// the expiry time for the presigned URL
	expiry := 5 * time.Minute
	if method == "PUT" {
		req, _ := c.client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectName),
		})

		url, err := req.Presign(expiry)
		if err != nil {
			log.Printf("Failed to presign PUT request for %s in bucket %s: %v", objectName, bucketName, err)
			return "", err
		}
		return url, nil

	} else if method == "GET" {
		req, _ := c.client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectName),
		})
		url, err := req.Presign(expiry)
		if err != nil {
			log.Printf("Failed to presign GET request for %s in bucket %s: %v", objectName, bucketName, err)
			return "", err
		}
		return url, nil
	}
	return "", nil
}
