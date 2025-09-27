package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {

	endpoint := "https://s3.regru.cloud"
	accessKey := "YS6GNF9VB0C96CDD63FG"
	secretKey := "kHGLbd7xtmMrNuoHuo249owr54CPvN8B5vVOEHc0"
	region := "ru-central1"

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	// Тестируем подключение
	file, err := os.Open("lol.txt")

	if err != nil {
		log.Fatal(err)
	}

	buc := "electronic"
	key := "product/lol.txt"
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &buc,
		Key:    &key,
		Body:   file,
	})

	if err != nil {
		log.Fatal(err)
	}

}
