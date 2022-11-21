package client

import (
	"net/http"

	"net/url"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

// ClientCreator creates a client for this service.
type ClientCreator func(url.URL, int) Client

type Client interface {
	Request(*http.Request) (*http.Response, error)
}

type client struct {
	address string
	index   int
	cli     *http.Client
	url     url.URL
}

func NewClient(address url.URL, index int) Client {
	return &client{
		address: address.String(),
		index:   index,
		cli:     cleanhttp.DefaultPooledClient(),
		url:     address,
	}
}

func (c *client) Request(req *http.Request) (*http.Response, error) {
	// Override the request with the right new url.
	req.URL = &url.URL{
		Scheme:   c.url.Scheme,
		Host:     c.url.Host,
		User:     c.url.User,
		Path:     req.URL.Path,
		RawPath:  req.URL.RawPath,
		RawQuery: req.URL.RawQuery,
		Opaque:   req.URL.Opaque,
		Fragment: req.URL.Fragment,
	}

	return c.cli.Do(req)
}
