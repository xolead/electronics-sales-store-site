package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3StorageConfig struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Region    string
	Bucket    string
	Folder    string
}

func main() {
	cfg := S3StorageConfig{
		AccessKey: "YS6GNF9VB0C96CDD63FG",
		SecretKey: "kHGLbd7xtmMrNuoHuo249owr54CPvN8B5vVOEHc0",
		Endpoint:  "https://s3.regru.cloud",
		Region:    "ru-central1",
		Bucket:    "electronic",
		Folder:    "products",
	}

	err := setupS3CORS(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to setup CORS: %v", err)
	}

	fmt.Println("✅ CORS configured successfully!")
}

func setupS3CORS(cfg S3StorageConfig) error {
	clientCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client := s3.NewFromConfig(clientCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
	})

	_, err = client.PutBucketCors(context.TODO(), &s3.PutBucketCorsInput{
		Bucket: aws.String(cfg.Bucket),
		CORSConfiguration: &types.CORSConfiguration{
			CORSRules: []types.CORSRule{
				{
					AllowedHeaders: []string{"*"},
					AllowedMethods: []string{"PUT", "POST", "GET", "DELETE", "HEAD"},
					AllowedOrigins: []string{
						"http://localhost:3000",
						"https://localhost:3000",
						"http://127.0.0.1:3000",
					},
					ExposeHeaders: []string{"ETag"},
					MaxAgeSeconds: aws.Int32(3000),
				},
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to setup CORS: %w", err)
	}

	fmt.Printf("✅ CORS configured for bucket: %s\n", cfg.Bucket)
	return nil
}
