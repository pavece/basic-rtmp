package uploader

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var FileUploaderInstance FileUploader



type FileUploader struct {
	s3Client *s3.Client
}

func initFileUploader() *s3.Client{
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("Error loading AWS config: ", err)
	}

	s3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		customEndpoint := os.Getenv("S3_ENDPOINT")
		if customEndpoint != "" {
			o.BaseEndpoint = aws.String(customEndpoint)
			o.UsePathStyle = true
		}
	})
	
	_, err = s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err == nil {
		fmt.Println("Connected to S3 / MinIO")
	}
	return s3Client
}

func (u *FileUploader) SetupFileUploader(){
	if u.s3Client == nil {
		u.s3Client = initFileUploader()
	}
}

func (u *FileUploader) UploadFile(fileReader io.Reader, destName string) error {
	FileUploaderInstance.SetupFileUploader()
	
	_, err := u.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("CDN_BUCKET_NAME")),
		Key:    aws.String(destName),
		Body:   fileReader,
	})
	
	if err != nil {
		return err
	}

	return nil
}