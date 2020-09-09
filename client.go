package snorlax

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// DefaultClient is a Snorlax client configured with all of the default options.
var DefaultClient = &Client{
	BaseURL: "",

	httpClient:   http.DefaultClient,
	requestHooks: make([]RequestHook, 0),
	proxyURL:     nil,
}

// Client defines a stateful REST client able to perform HTTP requests.
type Client struct {
	// BaseURL is prepended to the URI of all requests made by the Client.
	BaseURL string

	httpClient *http.Client
	proxyURL   *url.URL
	// requestHooks is a list of RequestHooks that get applied in order to
	// requests made by the Client just before they are sent.
	requestHooks []RequestHook
}

func (c *Client) call(ctx context.Context, method, target string,
	query url.Values, body io.Reader, hooks ...RequestHook) (*Response,
	error) {

	u := strings.Join([]string{c.BaseURL, target}, "")
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
	for _, hook := range append(c.requestHooks, hooks...) {
		hook(req)
	}

	// httpClient could be nil if it was not set on Client creation. This is to
	// ensure that a default http client is always present to prevent panics.
	if c.httpClient == nil {
		c.createClient()
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform http request: %w", err)
	}

	return &Response{*res}, nil
}

func (c *Client) createClient() {
	c.httpClient = &http.Client{}
}

// AddRequestHook adds a RequestHook to be run just before the client
// performs requests. These hooks are executed in order.
func (c *Client) AddRequestHook(hook RequestHook) *Client {
	if c.requestHooks == nil {
		c.requestHooks = make([]RequestHook, 0)
	}

	c.requestHooks = append(c.requestHooks, hook)
	return c
}

// AddRequestHooks adds multiple RequestHooks to be run just before the client
// performs requests. These hooks are executed in order.
func (c *Client) AddRequestHooks(hooks ...RequestHook) *Client {
	for _, hook := range hooks {
		c.AddRequestHook(hook)
	}

	return c
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodDelete, target, query, body, opts...)
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodGet, target, query, nil, opts...)
}

// Head performs a HEAD request.
func (c *Client) Head(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodHead, target, query, nil, opts...)
}

// Options performs a OPTIONS request.
func (c *Client) Options(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodOptions, target, query, nil, opts...)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodPost, target, query, body, opts...)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodPut, target, query, body, opts...)
}

// RemoveProxy removes the currently set proxy.
func (c *Client) RemoveProxy() *Client {
	t, ok := c.httpClient.Transport.(*http.Transport)
	if !ok {
		// TODO: Add logging as an indication that this was skipped.
		return c
	}

	c.proxyURL = nil
	t.Proxy = nil

	return c
}

// SetBaseURL sets the url that is prepended to all request URLs.
func (c *Client) SetBaseURL(u string) *Client {
	if _, err := url.Parse(u); err != nil {
		// TODO: Add logging as an indication that this failed.
		return c
	}

	c.BaseURL = u
	return c
}

// SetRequestHooks sets the RequestHooks to be run just before the client
// performs requests. These are run in order. Calling SetRequestHooks will
// replace any existing RequestHooks that have been added. To add RequestHooks
// without replacing other hooks use AddRequestHook(s).
func (c *Client) SetRequestHooks(hooks []RequestHook) {
	c.requestHooks = hooks
}

// SetProxy sets the proxy URL for the Snorlax client. If the provided URL fails
// to be parsed then nothing will be set.
func (c *Client) SetProxy(u string) *Client {
	t, ok := c.httpClient.Transport.(*http.Transport)
	if !ok {
		// TODO: Add logging as an indication that this failed.
		return c
	}

	proxyURL, err := url.Parse(u)
	if err != nil {
		// TODO: Add logging as an indication that this failed.
		return c
	}

	c.proxyURL = proxyURL
	t.Proxy = http.ProxyURL(proxyURL)

	return c
}
