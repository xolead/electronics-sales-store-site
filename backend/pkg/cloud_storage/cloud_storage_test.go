package cloudstorage_test

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	cloudstorage "electronic/pkg/cloud_storage"
)

var cfg cloudstorage.S3StorageConfig = cloudstorage.S3StorageConfig{
	AccessKey: "YS6GNF9VB0C96CDD63FG",
	SecretKey: "kHGLbd7xtmMrNuoHuo249owr54CPvN8B5vVOEHc0",
	Endpoint:  "https://s3.regru.cloud",
	Region:    "ru-central1",
	Bucket:    "electronic",
	Folder:    "tests",
}

func TestUploadDonloadDeleteURL(t *testing.T) {

	//Загрузка
	s3s, err := cloudstorage.NewS3S(cfg)
	assert.NoError(t, err)

	name := "omg.txt"

	image, err := s3s.UploadURL(name)
	assert.NoError(t, err)

	file, err := os.Open("tests/" + name)
	assert.NoError(t, err)
	defer file.Close()

	key := image.FileID
	req, err := http.NewRequest("PUT", image.URL, file)
	assert.NoError(t, err)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// Скачивание

	image, err = s3s.DownloadURL(key)

	resp, err = http.Get(image.URL)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	downloadName := "omg_download.txt"
	file, err = os.Create("tests/" + downloadName)
	assert.NoError(t, err)

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	assert.NoError(t, err)

	// Удаление
	image, err = s3s.DeleteURL(key)
	assert.NoError(t, err)

	req, err = http.NewRequest("DELETE", image.URL, nil)

	assert.NoError(t, err)

	resp, err = client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
}
