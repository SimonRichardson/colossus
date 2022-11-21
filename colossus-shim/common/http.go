package common

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func newRequest(reqType string, url string, payload []byte) (*http.Request, error) {
	if payload == nil {
		return http.NewRequest(reqType, url, nil)
	}
	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.ContentLength = int64(len(payload))
	return req, err
}

func request(reqType string, url string, payload []byte, fn func(http.Header)) ([]byte, error) {
	req, err := newRequest(reqType, url, payload)
	if err != nil {
		return nil, err
	}

	fn(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if errored(resp.StatusCode) {
		return nil, fmt.Errorf("Request error: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Get(url string, fn func(http.Header)) ([]byte, error) {
	return request("GET", url, nil, fn)
}

func Post(url string, bytes []byte, fn func(http.Header)) ([]byte, error) {
	return request("POST", url, bytes, fn)
}

func errored(s int) bool {
	return !(s == http.StatusOK || s == http.StatusCreated || s == http.StatusNoContent)
}
