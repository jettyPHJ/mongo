package myhttp_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fmt.Printf("请求参数：%+v", params)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)

	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func Test001(t *testing.T) {
	extraParams := map[string]string{
		"Meta": `{"factoryID":"factoryID_002","startTime":1666920172835,"endTime":1666920244224,"traceFileName":"0a4b178d-4b66-4ad4-8f9c-fbba26441277.csv","objectClass":"car","roadNumber":"road001","traceID":"0a4b178d-4b66-4ad4-8f9c-fbba26441277","licensePlate":"cd123455","traceNumber":3}`,
	}
	request, err := newfileUploadRequest("http://localhost:8080/trajectoryFile", extraParams, "file", "/home/jette/cloud_trackfile/testSimulator/tem/0a4b178d-4b66-4ad4-8f9c-fbba26441277.csv")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
	}
}
