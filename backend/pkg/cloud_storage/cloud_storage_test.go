// cloud_storage_test.go
package cloudstorage_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	cloudstorage "electronic/pkg/cloud_storage"
)

type CloudStorageTestSuite struct {
	suite.Suite
	storage cloudstorage.CloudStorage
	cfg     cloudstorage.S3StorageConfig
}

func (suite *CloudStorageTestSuite) SetupSuite() {
	suite.cfg = cloudstorage.S3StorageConfig{
		AccessKey: "YS6GNF9VB0C96CDD63FG",
		SecretKey: "kHGLbd7xtmMrNuoHuo249owr54CPvN8B5vVOEHc0",
		Endpoint:  "https://s3.regru.cloud",
		Region:    "ru-central1",
		Bucket:    "electronic",
		Folder:    "tests",
	}

	var err error
	suite.storage, err = cloudstorage.NewS3S(suite.cfg)
	require.NoError(suite.T(), err, "Failed to initialize cloud storage")
}

func (suite *CloudStorageTestSuite) TestUploadDownloadDelete() {
	t := suite.T()

	// Test data
	testFileName := "test_file.txt"
	testFilePath := filepath.Join("testdata", testFileName)
	downloadFileName := "test_file_download.txt"
	downloadFilePath := filepath.Join("testdata", downloadFileName)

	// Create test directory and file
	suite.createTestFile(testFilePath, "test file content")
	defer suite.cleanupTestFiles(testFilePath, downloadFilePath)

	// Upload test
	t.Run("Upload File", func(t *testing.T) {
		uploadResult, err := suite.storage.UploadURL(testFileName)
		assert.NoError(t, err)
		assert.NotEmpty(t, uploadResult.URL)
		assert.NotEmpty(t, uploadResult.FileID)

		err = suite.executeUpload(uploadResult.URL, testFilePath)
		assert.NoError(t, err)
	})

	// Download test
	t.Run("Download File", func(t *testing.T) {
		downloadResult, err := suite.storage.DownloadURL(testFileName)
		assert.NoError(t, err)
		assert.NotEmpty(t, downloadResult.URL)

		err = suite.executeDownload(downloadResult.URL, downloadFilePath)
		assert.NoError(t, err)

		// Verify file content
		originalContent, _ := os.ReadFile(testFilePath)
		downloadedContent, _ := os.ReadFile(downloadFilePath)
		assert.Equal(t, originalContent, downloadedContent)
	})

	// Delete test
	t.Run("Delete File", func(t *testing.T) {
		deleteResult, err := suite.storage.DeleteURL(testFileName)
		assert.NoError(t, err)
		assert.NotEmpty(t, deleteResult.URL)

		err = suite.executeDelete(deleteResult.URL)
		assert.NoError(t, err)
	})
}

func (suite *CloudStorageTestSuite) createTestFile(filePath, content string) {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	require.NoError(suite.T(), err)

	err = os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(suite.T(), err)
}

func (suite *CloudStorageTestSuite) cleanupTestFiles(files ...string) {
	for _, file := range files {
		os.Remove(file)
	}
}

func (suite *CloudStorageTestSuite) executeUpload(uploadURL, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %d", resp.StatusCode)
	}
	return nil
}

func (suite *CloudStorageTestSuite) executeDownload(downloadURL, filePath string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (suite *CloudStorageTestSuite) executeDelete(deleteURL string) error {
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete failed with status: %d", resp.StatusCode)
	}
	return nil
}

func TestCloudStorageTestSuite(t *testing.T) {
	suite.Run(t, new(CloudStorageTestSuite))
}
