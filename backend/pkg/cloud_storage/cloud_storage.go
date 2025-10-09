package cloudstorage

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"electronic/pkg/models"
)

type CloudStorage interface {
	UploadURL(filename string) (models.S3SImage, error)
	DownloadURL(key string) (models.S3SImage, error)
	DeleteURL(key string) (models.S3SImage, error)
}

type s3Storage struct {
	*s3.Client
	bucket string
	folder string
}

type S3StorageConfig struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Region    string
	Bucket    string
	Folder    string
}

func LoadConfig() S3StorageConfig {
	return S3StorageConfig{
		AccessKey: os.Getenv("access_key"),
		SecretKey: os.Getenv("secret_key"),
		Endpoint:  os.Getenv("endpoint"),
		Region:    os.Getenv("region"),
		Bucket:    os.Getenv("bucket"),
		Folder:    os.Getenv("folder"),
	}
}

func NewS3S(cfg S3StorageConfig) (CloudStorage, error) {
	log.Println("S3S: ", cfg)
	clientCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("NewS3S ошибка создание clientCfg: %w", err)
	}

	clientS3S := s3.NewFromConfig(clientCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
	})

	s3s := &s3Storage{
		Client: clientS3S,
		bucket: cfg.Bucket,
		folder: cfg.Folder,
	}
	return s3s, nil
}

func (s3s *s3Storage) UploadURL(filename string) (models.S3SImage, error) {
	presignClient := s3.NewPresignClient(s3s.Client)

	mime := filepath.Ext(filename)
	fileID := uuid.New().String() + mime
	image := models.S3SImage{FileID: fileID}

	expires := 10 * time.Minute

	presignResult, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s3s.bucket),
		Key:         aws.String(s3s.folder + "/" + fileID),
		ContentType: aws.String(mime),
	}, s3.WithPresignExpires(expires))

	if err != nil {
		return image, fmt.Errorf("UploadURL ошибка созданяи предподписаного url: %w", err)
	}

	image.URL = presignResult.URL

	return image, nil

}

func (s3s *s3Storage) DownloadURL(key string) (models.S3SImage, error) {
	presignClient := s3.NewPresignClient(s3s.Client)

	image := models.S3SImage{
		FileID: key,
	}
	expires := 30 * time.Minute

	presignResult, err := presignClient.PresignGetObject(
		context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(s3s.bucket),
			Key:    aws.String(s3s.folder + "/" + key)},
		s3.WithPresignExpires(expires),
	)

	if err != nil {
		return image, fmt.Errorf("DownloadURL ошибка создания предподписанного url: %w", err)
	}

	image.URL = presignResult.URL
	return image, nil

}

func (s3s *s3Storage) DeleteURL(key string) (models.S3SImage, error) {
	image := models.S3SImage{
		FileID: key,
	}
	presignClient := s3.NewPresignClient(s3s.Client)

	presignResult, err := presignClient.PresignDeleteObject(
		context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(s3s.bucket),
			Key:    aws.String(s3s.folder + "/" + key)},
		s3.WithPresignExpires(15*time.Minute),
	)

	if err != nil {
		return image, fmt.Errorf("DeleteURL ошибка создания предподписаного url: %w", err)
	}

	image.URL = presignResult.URL
	return image, nil

}
