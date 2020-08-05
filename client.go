package transit

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type client struct {
	opts clientOptions
}

// NewClient returns a transit Client configured with the provided
// ClientOptions.
func NewClient(opts ...ClientOption) *client {
	c := client{
		opts: defaultOptions,
	}

	for _, o := range opts {
		o(&c.opts)
	}

	return &c
}

func (c *client) Delete(uri string, params url.Values, body io.Reader) (
	*http.Response, error) {
	req, err := makeRequest(http.MethodDelete, uri, params, body)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *client) Get(uri string, params url.Values) (*http.Response, error) {
	req, err := makeRequest(http.MethodGet, uri, params, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *client) Post(uri string, params url.Values, body io.Reader) (
	*http.Response, error) {
	req, err := makeRequest(http.MethodPost, uri, params, body)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *client) Put(uri string, params url.Values, body io.Reader) (
	*http.Response, error) {
	req, err := makeRequest(http.MethodPut, uri, params, body)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func makeRequest(method, uri string, params url.Values, body io.Reader) (
	*http.Request, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	u.RawQuery = params.Encode()
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *client) Do(r *http.Request) (*http.Response, error) {
	// Set the request headers.
	for k, v := range c.opts.headers {
		r.Header.Set(k, v[0])
	}

	// Re-parse the request URI if a base URL has been set in the client.
	if c.opts.baseURL != "" {
		newURL := fmt.Sprintf("%s%s", c.opts.baseURL, r.URL.String())
		u, err := url.Parse(newURL)
		if err != nil {
			return nil, err
		}
		r.URL = u
	}

	// Execute all request adapters.
	for _, adapter := range c.opts.adapters {
		adapter.Adapt(r)
	}

	res, err := c.opts.transport.Do(r)
	if err != nil {
		return nil, err
	}

	return res, nil
}


