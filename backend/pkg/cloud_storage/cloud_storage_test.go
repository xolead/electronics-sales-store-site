package cloudstorage_test

import (
	"bytes"
	cloudstorage "electronic/pkg/cloud_storage"
	"io"
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

	fileId := UploadURL(s3s, tests, t)
	DownloadURL(s3s, tests, fileId, t)
	DeleteURL(s3s, tests, fileId, t)

	assert.NoError(t, err)

}

func UploadURL(s3s cloudstorage.CloudStorage, tests []test, t *testing.T) []string {
	fileId := make([]string, 0, len(tests))
	for _, test := range tests {
		image, err := s3s.UploadURL(test.fileName)
		assert.NoError(t, err)
		fileId = append(fileId, image.FileID)

		file, err := os.Create("tests/" + test.fileName)
		assert.NoError(t, err)

		_, err = file.WriteString(test.data)
		assert.NoError(t, err)
		err = file.Close()
		assert.NoError(t, err)

		fileData, err := os.ReadFile("tests/" + test.fileName)
		assert.NoError(t, err)

		client := &http.Client{}
		req, err := http.NewRequest(
			"PUT",
			image.URL,
			bytes.NewReader(fileData),
		)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

	}

	return fileId
}

func DownloadURL(s3s cloudstorage.CloudStorage, tests []test, fileId []string, t *testing.T) {

	for i, test := range tests {
		image, err := s3s.DownloadURL(fileId[i])
		assert.NoError(t, err)

		resp, err := http.Get(image.URL)
		assert.NoError(t, err)
		defer resp.Body.Close()

		downloadedData, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, test.data, string(downloadedData))
	}
}

func DeleteURL(s3s cloudstorage.CloudStorage, tests []test, fileId []string, t *testing.T) {
	for i, _ := range tests {
		image, err := s3s.DeleteURL(fileId[i])
		assert.NoError(t, err)

		client := &http.Client{}
		req, err := http.NewRequest("DELETE", image.URL, nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)

		resp.Body.Close()
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent)

		resp, err = http.Get(image.URL)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.NotEqual(t, http.StatusOK, resp.StatusCode)

	}
}
