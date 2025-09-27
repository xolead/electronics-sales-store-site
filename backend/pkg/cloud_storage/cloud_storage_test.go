package cloudstorage_test

import (
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

func TestUploadURL(t *testing.T) {
	s3s, err := cloudstorage.NewS3S(cfg)
	assert.NoError(t, err)

	image, err := s3s.UploadURL("leva.txt")
	assert.NoError(t, err)

	file, err := os.Open("test.txt")
	assert.NoError(t, err)
	defer file.Close()

	req, err := http.NewRequest("PUT", image.URL, file)
	assert.NoError(t, err)

	client := &http.Client{}

	resp, err := client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
