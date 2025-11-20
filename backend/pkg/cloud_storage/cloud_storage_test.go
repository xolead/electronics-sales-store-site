package cloudstorage_test

import (
	"bytes"
	cloudstorage "electronic/pkg/cloud_storage"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cfg cloudstorage.S3StorageConfig = cloudstorage.S3StorageConfig{
	AccessKey: "YS6GNF9VB0C96CDD63FG",
	SecretKey: "kHGLbd7xtmMrNuoHuo249owr54CPvN8B5vVOEHc0",
	Endpoint:  "https://s3.regru.cloud",
	Region:    "ru-central1",
	Bucket:    "electronic",
	Folder:    "tests",
}

type test struct {
	fileName string
	data     string
}

func TestCloudStorage(t *testing.T) {
	s3s, err := cloudstorage.NewS3S(cfg)

	var tests = []test{
		{"test1.txt", "lolsasdfadwdqwwd"},
		{"test2.txt", "123fdfqsd"},
		{"test3.txt", ">>/"},
	}

	assert.NoError(t, err)

}

func UploadURL(s3s cloudstorage.CloudStorage, tests []test, t *testing.T) {

	for _, test := range tests {
		image, err := s3s.UploadURL(test.fileName)
		assert.NoError(t, err)

		file, err := os.Create("test/" + test.fileName)
		assert.NoError(t, err)

		_, err = file.WriteString(test.data)
		assert.NoError(t, err)

		client := &http.Client{}
		req, err := http.NewRequest("PUT", presignedURL, bytes.NewReader(data))
		assert.NoError(t, err)

		err = file.Close()
		assert.NoError(t, err)

	}

}
