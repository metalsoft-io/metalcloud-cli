package networking

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
	RequestTypeGet = "GET"
	RequestTypePost = "POST"

	CatalogUrlPath = "/ms-api/firmware/catalog"
	BinaryUrlPath  = "/ms-api/firmware/binary"
)

func SendMsRequest(requestType, url, apiKey string, jsonData []byte) (string, error) {
	req, err := createHttpRequest(requestType, url, apiKey, jsonData)
	if err != nil {
		return "", err
	}

	setHttpRequestHeaders(req, requestType, apiKey)
	body, err := sendHttpRequest(req)
	if err != nil {
		return "", err
	}

	return body, nil
}

func createHttpRequest(requestType, url, apiKey string, jsonData []byte) (*http.Request, error) {
	var req *http.Request
	var err error

	switch requestType {
	case RequestTypeGet:
		req, err = http.NewRequest(requestType, url, nil)
		if err != nil {
			return req, err
		}
	case RequestTypePost:
		req, err = http.NewRequest(requestType, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return req, err
		}
	default:
		supportedRequestTypes := []string{RequestTypeGet, RequestTypePost}
		return req, fmt.Errorf("unsupported request type %s. Supported request types are %v", requestType, supportedRequestTypes)
	}

	return req, nil
}

func setHttpRequestHeaders(req *http.Request, requestType, apiKey string) {
	var bearer = "Bearer " + apiKey
	req.Header.Add("Authorization", bearer)

	if requestType == RequestTypePost {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}
}

func sendHttpRequest(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}
