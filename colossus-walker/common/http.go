package common

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/SimonRichardson/colossus/typex"
)

func NewRequest(reqType string, url string, payload []byte) (*http.Request, error) {
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

func Request(reqType string, url string, payload []byte) ([]byte, error) {
	req, err := NewRequest(reqType, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if errored(resp.StatusCode) {
		typex.Fatal(fmt.Errorf("Request error: %d", resp.StatusCode))
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

func Get(url string) ([]byte, error) {
	return Request("GET", url, nil)
}

func errored(s int) bool {
	return !(s == http.StatusOK || s == http.StatusCreated || s == http.StatusNoContent)
}
