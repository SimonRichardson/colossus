package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

type ClientCreator func(string, string, *ConnectionTimeout) Client

type Options struct {
	Versioned bool
}

var (
	Versioned   = Options{true}
	Unversioned = Options{false}
)

type Response struct {
	Status     string
	StatusCode int
	Proto      string

	Header http.Header
	Bytes  []byte
}

type Client interface {
	Get(string, Options, func(http.Header)) (*Response, error)
	Post(string, []byte, Options, func(http.Header)) (*Response, error)
}

type client struct {
	address, version string
	client           *http.Client
}

func New(address, version string, timeout *ConnectionTimeout) Client {
	return &client{
		address,
		version,
		cleanhttp.DefaultPooledClient(),
	}
}

func (c *client) Get(url string, options Options, fn func(http.Header)) (*Response, error) {
	return c.request("GET", url, nil, options, fn)
}

func (c *client) Post(url string, payload []byte, options Options, fn func(http.Header)) (*Response, error) {
	return c.request("POST", url, payload, options, fn)
}

func (c *client) request(reqType string,
	url string,
	payload []byte,
	options Options,
	fn func(http.Header),
) (*Response, error) {
	address := c.address
	if options.Versioned {
		address = fmt.Sprintf("%s%s", address, c.version)
	}

	req, err := newRequest(reqType, fmt.Sprintf("%s%s", address, url), payload)
	if err != nil {
		return nil, err
	}

	fn(req.Header)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Make sure we drain the response.
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	return readResponse(resp)
}

func newRequest(reqType string, url string, payload []byte) (req *http.Request, err error) {
	if payload == nil {
		req, err = http.NewRequest(reqType, url, nil)
	} else {
		req, err = http.NewRequest(reqType, url, bytes.NewBuffer(payload))
	}
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return
}

func readResponse(resp *http.Response) (*Response, error) {
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Proto:      resp.Proto,

		Header: resp.Header,
		Bytes:  bytes,
	}, nil
}
