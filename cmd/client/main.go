package client

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type TUser struct {
	C http.Client
}

func NewTUser(timeOut int, needRedirect bool) *TUser {
	// Using or not using Redirect
	var fR func(*http.Request, []*http.Request) error
	if !needRedirect {
		fR = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		fR = func(req *http.Request, via []*http.Request) error {
			return nil
		}
	}

	// Add Cookiejar
	jar, _ := cookiejar.New(nil)

	return &TUser{
		C: http.Client{
			Timeout:       time.Duration(time.Duration(timeOut) * time.Second),
			CheckRedirect: fR,
			Jar:           jar,
		},
	}
}

/*
ReadBodyDefault - чтобы TCP-соединение использовалось повторно, клиент должен
обязательно прочитать тело ответа до конца и закрыть
*/
func (tu *TUser) ReadBodyDefault(resBody io.ReadCloser) error {
	_, err := io.Copy(io.Discard, resBody)
	if err != nil {
		return err
	}
	resBody.Close()
	return nil
}

func (tu *TUser) PreparSendingFile(nameFile string) (body *bytes.Buffer, content string, err error) {
	// Open file
	file, err := os.Open(nameFile)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	// Create buffer
	body = &bytes.Buffer{}
	// Create multipart.Writer
	writer := multipart.NewWriter(body)
	// prepar special form
	part, err := writer.CreateFormFile("uploadfile", nameFile)
	if err != nil {
		return nil, "", err
	}

	// Copy file to form
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, "", err
	}
	writer.Close()
	return body, writer.FormDataContentType(), nil
}

func (tu *TUser) SendWithFile(url string, nameFile string) (response *http.Response, err error) {
	body, content, err := tu.PreparSendingFile(nameFile)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	// Add Header
	req.Header.Set("Content-Type", content)
	return tu.C.Do(req)
}

func (tu *TUser) SendWithQuery(method, uurl string, query ...string) (response *http.Response, err error) {
	if len(query)&1 != 0 {
		return nil, errors.New("number of elements must be even")
	}
	if len(query) < 2 {
		return nil, errors.New("number of elements must be at least 2")
	}

	data := url.Values{}
	value := 1
	key := 0

	// Installation query
	for value < len(query) {
		data.Set(query[key], query[value])
		value++
		key++
	}

	// Write request
	req, err := http.NewRequest(method, uurl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Add Headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	return tu.C.Do(req)
}
