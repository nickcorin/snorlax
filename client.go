package snorlax

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client defines a stateful REST client able to perform HTTP requests.
type Client struct {
	opts *ClientOptions
}

// New returns a snorlax client configured with the provided
// ClientOptions.
func New(opts *ClientOptions) *Client {
	if opts == nil {
		opts = &defaultOptions
	}

	if opts.Transport == nil {
		opts.Transport = defaultOptions.Transport
	}

	if opts.CallOptions == nil {
		opts.CallOptions = defaultOptions.CallOptions
	}

	return &Client{opts}
}

func (c *Client) call(ctx context.Context, method, path string,
	query url.Values, body io.Reader, opts ...CallOption) (*Response,
	error) {

	u := strings.Join([]string{c.opts.BaseURL, path}, "")
	uri, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %s: %w", uri, err)
	}

	// Once the client has a logger built in, this can be logged as a warning
	// rather than returned as an error.
	if uri.RawQuery != "" {
		return nil, fmt.Errorf("query params should not be set on the path")
	}
	uri.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, uri.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	// We first apply the request options from the client, so that they can be
	// optionally overridden by individual request options.
	for _, opt := range append(c.opts.CallOptions, opts...) {
		opt.Apply(req)
	}

	res, err := c.opts.Transport.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform http request: %w", err)
	}

	return &Response{*res}, nil
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, uri string, query url.Values,
	body io.Reader, opts ...CallOption) (*Response, error) {
	return c.call(ctx, http.MethodDelete, uri, query, body, opts...)
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, uri string, query url.Values,
	opts ...CallOption) (*Response, error) {
	return c.call(ctx, http.MethodGet, uri, query, nil, opts...)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, uri string, query url.Values,
	body io.Reader, opts ...CallOption) (*Response, error) {
	return c.call(ctx, http.MethodPost, uri, query, body, opts...)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, uri string, query url.Values,
	body io.Reader, opts ...CallOption) (*Response, error) {
	return c.call(ctx, http.MethodPut, uri, query, body, opts...)
}
