package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
	Endpoint  string
}

type S3Client struct {
	client   *s3.Client
	uploader *manager.Uploader
	bucket   string
	region   string
	baseURL  string
}

func main() {
	cfg := S3Config{
		AccessKey: os.Getenv("YS6GNF9VB0C96CDD63FG"),
		SecretKey: os.Getenv("kHGLbd7xtmMrNuoHuo249owr54CPvN8B5vVOEHc0"),
		Region:    os.Getenv("ru"),
		Bucket:    os.Getenv("electronic"),
		Endpoint:  os.Getenv("https://s3.regru.cloud/"),
	}
	creds := credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "wkefjwo")

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds))
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(awsCfg)

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024
		u.Concurrency = 5
	})

	baseUrl := "https://s3.regru.cloud/electronic/"

	s := S3Client{
		client:   client,
		uploader: uploader,
		bucket:   cfg.Bucket,
		baseURL:  baseUrl,
	}

	file, err := os.Open("lol.txt")
	if err != nil {
		log.Fatal(err)
	}

	result, err := s.uploader.Upload(
		context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(s.bucket),
			Key:         aws.String("lol.txt"),
			Body:        file,
			ContentType: aws.String("txt"),
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
