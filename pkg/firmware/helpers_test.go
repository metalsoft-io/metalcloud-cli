package firmware

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jlaffaye/ftp"
)

func TestUploadFileToRepo(t *testing.T) {
	ftpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "220 FTP Server Ready")
	}))
	defer ftpServer.Close()

	mockBinaryData := []byte("binary file")
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(mockBinaryData)
	}))
	defer mockServer.Close()

	ftpServerURL := ftpServer.URL[7:] // remove "http://"
	repoHost := ftpServerURL
	repoPath := "/path/to/repo"
	repoUser := "username"
	repoPassword := "password"

	url := mockServer.URL
	fileName := "mock_file.bin"

	testFTPClient, err := ftp.Connect(ftpServerURL)
	if err != nil {
		t.Fatalf("Failed to connect to the test FTP server: %s", err)
	}
	defer testFTPClient.Quit()

	err = testFTPClient.Login(repoUser, repoPassword)
	if err != nil {
		t.Fatalf("Failed to log in to the test FTP server: %s", err)
	}

	err = testFTPClient.ChangeDir(repoPath)
	if err != nil {
		t.Fatalf("Failed to change directory in the test FTP server: %s", err)
	}

	err = uploadFileToRepo(url, fileName, repoHost, repoPath, repoUser, repoPassword)
	if err != nil {
		t.Errorf("uploadFileToRepo returned an error: %s", err)
	}

	file, err := testFTPClient.Retr(fileName)
	if err != nil {
		t.Fatalf("Failed to retrieve the uploaded file from the test FTP server: %s", err)
	}
	defer file.Close()

	uploadedData, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read the uploaded file data from the test FTP server: %s", err)
	}

	if string(uploadedData) != string(mockBinaryData) {
		t.Errorf("Uploaded file data does not match the mock binary data")
	}
}

func TestMain(m *testing.M) {
	m.Run()
}