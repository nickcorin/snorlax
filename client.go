package snorlax

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DefaultClient is a Snorlax client configured with all of the default options.
var DefaultClient = &Client{
	BaseURL:       "",
	EnableMetrics: false,

	headers:      make(http.Header),
	httpClient:   http.DefaultClient,
	requestHooks: make([]RequestHook, 0),
	proxyURL:     nil,
}

// Client is a stateful REST client that is able to make HTTP requests.
type Client struct {
	// BaseURL is prepended to the URI of all requests made by the Client.
	BaseURL string
	// EnableMetrics enables prometheus metrics. This is disabled by default.
	EnableMetrics bool

	headers    http.Header
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

	// TODO: Replace this error with a warning log once a logger has been added
	// to the client. We shouldn't add logs until there is a configuration to
	// disable them.
	if uri.RawQuery != "" {
		return nil, fmt.Errorf("query params should not be set on the path")
	}
	uri.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, uri.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	// Set the request headers with all the headers configured in the client.
	if c.headers == nil {
		c.headers = make(http.Header)
	}
	req.Header = c.headers

	// Automatically add the Content-Length header.
	req.Header.Set(http.CanonicalHeaderKey("Content-Length"),
		strconv.FormatInt(req.ContentLength, 10))

	// We first apply the request options from the client, so that they can be
	// optionally overridden by individual request options.
	for _, hook := range append(c.requestHooks, hooks...) {
		hook(req)
	}

	// httpClient is usually nil on the first request made by the client. This
	// prevents panics by using the http.DefaultClient. In most cases, this will
	// be sufficient. In cases where the caller wants more control over the
	// client's configuration - SetHTTPClient can be used.
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	reqStart := time.Now()
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform http request: %w", err)
	}

	// TODO: Rethink how to provide the target path as a label effectively.
	// Clients sending requests to dynamic paths can overload prometheus.
	if c.EnableMetrics {
		latencyHist.WithLabelValues(method,
			strconv.Itoa(res.StatusCode), req.URL.Path).Observe(
			time.Since(reqStart).Seconds())
	}

	return &Response{*res}, nil
}

// AddHeader appends a header value to the Client to be sent in every request.
// To replace the current existing header use SetHeader.
func (c *Client) AddHeader(key, value string) *Client {
	c.headers.Add(key, value)
	return c
}

// AddRequestHook appends a RequestHook to the list of hooks which are to be run
// just before the client sends a request. RequestHooks are executed in the
// order they are added.
func (c *Client) AddRequestHook(hook RequestHook) *Client {
	if c.requestHooks == nil {
		c.requestHooks = make([]RequestHook, 0)
	}

	c.requestHooks = append(c.requestHooks, hook)
	return c
}

// AddRequestHooks is a convenience function which calls AddRequestHook multiple
// times.
func (c *Client) AddRequestHooks(hooks ...RequestHook) *Client {
	for _, hook := range hooks {
		c.AddRequestHook(hook)
	}

	return c
}

// Delete performs a delete request using the DefaultClient.
//
// You can optionally configure the request using RequestHooks. If you require
// finer control to configure every request made then construct a Client and
// configure it as needed.
func Delete(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodDelete, target, query, body,
		opts...)
}

// Delete performs a delete request.
//
// You can optionally configure the request using RequestHooks, or by confiuring
// the client if you need to configure all requests.
func (c *Client) Delete(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodDelete, target, query, body, opts...)
}

// Get performs a get request using the DefaultClient.
//
// You can optionally configure the request using RequestHooks. If you require
// finer control to configure every request made then construct a Client and
// configure it as needed.
func Get(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodGet, target, query, nil, opts...)
}

// Get performs a Get request.
//
// You can optionally configure the request using RequestHooks, or by confiuring
// the client if you need to configure all requests.
func (c *Client) Get(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodGet, target, query, nil, opts...)
}

// Head performs a head request using the DefaultClient.
//
// You can optionally configure the request using RequestHooks. If you require
// finer control to configure every request made then construct a Client and
// configure it as needed.
func (c *Client) Head(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodHead, target, query, nil, opts...)
}

// Head performs a head request.
//
// You can optionally configure the request using RequestHooks, or by confiuring
// the client if you need to configure all requests.
func Head(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodHead, target, query, nil, opts...)
}

// Options performs an options request using the DefaultClient.
//
// You can optionally configure the request using RequestHooks. If you require
// finer control to configure every request made then construct a Client and
// configure it as needed.
func Options(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodOptions, target, query, nil,
		opts...)
}

// Options performs an options request.
//
// You can optionally configure the request using RequestHooks, or by confiuring
// the client if you need to configure all requests.
func (c *Client) Options(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodOptions, target, query, nil, opts...)
}

// Post performs a post request using the DefaultClient.
//
// You can optionally configure the request using RequestHooks. If you require
// finer control to configure every request made then construct a Client and
// configure it as needed.
func Post(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodPost, target, query, body,
		opts...)
}

// Post performs a post request.
//
// You can optionally configure the request using RequestHooks, or by confiuring
// the client if you need to configure all requests.
func (c *Client) Post(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodPost, target, query, body, opts...)
}

// Put performs a put request using the DefaultClient.
//
// You can optionally configure the request using RequestHooks. If you require
// finer control to configure every request made then construct a Client and
// configure it as needed.
func Put(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodPut, target, query, body, opts...)
}

// Put performs a put request.
//
// You can optionally configure the request using RequestHooks, or by confiuring
// the client if you need to configure all requests.
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

// SetHeader sets a header value in the client to be sent in every request. This
// will overwrite any exiting headers present associated with the same key. To
// add headers to the key instead of replacing them use AddHeader.
func (c *Client) SetHeader(key, value string) *Client {
	c.headers.Set(key, value)
	return c
}

// SetHTTPClient sets the internal http.Client that Snorlax uses to perform
// requests. Use this if you want to configure client internals like timeouts.
func (c *Client) SetHTTPClient(client *http.Client) *Client {
	c.httpClient = client
	return c
}

// SetRequestHooks sets the RequestHooks to be run just before the client
// performs requests. These are run in order. Calling SetRequestHooks will
// replace any existing RequestHooks that have been added. To add RequestHooks
// without replacing other hooks use AddRequestHook(s).
func (c *Client) SetRequestHooks(hooks []RequestHook) *Client {
	c.requestHooks = hooks
	return c
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
