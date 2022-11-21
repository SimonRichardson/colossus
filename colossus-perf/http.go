package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/tests"
	"github.com/SimonRichardson/colossus/typex"
	flatbuffers "github.com/google/flatbuffers/go"
)

func request(reqType string, url string, payload []byte) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(payload))
	if err != nil {
		typex.Fatal(err)
	}

	req.ContentLength = int64(len(payload))
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := client.Do(req)
	if err != nil {
		typex.Fatal(err)
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
		typex.Fatal(err)
	}

	return body
}

func errored(s int) bool {
	return !(s == http.StatusOK || s == http.StatusCreated || s == http.StatusNoContent)
}

func post(url string, payload []byte) []byte {
	return request("POST", url, payload)
}

func api(url string,
	key bson.ObjectId,
	maxSize int,
	co *coordinator.Coordinator,
) func(tests.PostBody) int {
	return func(values tests.PostBody) int {
		record := records.PostRecords{
			Records: values,
			MaxSize: int64(maxSize),
			Expiry:  defaultExpiry,
		}
		bytes, err := record.Write(flatbuffers.NewBuilder(0))
		if err != nil {
			typex.Fatal(err)
		}

		body := post(fmt.Sprintf("%s/http/v1/%s", url, key.Hex()), bytes)

		return len(body)
	}
}
