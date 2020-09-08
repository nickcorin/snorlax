package snorlax

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// DefaultClient is a snorlax client configured with all of the default options.
var DefaultClient = &client{
	opts: defaultOptions,
}

var defaultOptions = ClientOptions{
	PreRequestHooks: make([]RequestHook, 0),
	Transport:       http.DefaultClient,
}

// Client describes an HTTP REST client.
type Client interface {
	// Delete performs a DELETE request.
	Delete(context.Context, string, url.Values, io.Reader, ...RequestHook) (
		*Response, error)

	// Get performs a GET request.
	Get(context.Context, string, url.Values, ...RequestHook) (*Response, error)

	// Head performs a HEAD request.
	Head(context.Context, string, url.Values, ...RequestHook) (*Response, error)

	// Options performs an OPTIONS request.
	Options(context.Context, string, url.Values, ...RequestHook) (*Response,
		error)

	// Post performs a POST request.
	Post(context.Context, string, url.Values, io.Reader, ...RequestHook) (
		*Response, error)

	// Put performs a PUT request.
	Put(context.Context, string, url.Values, io.Reader, ...RequestHook) (
		*Response, error)
}

// ClientOptions defines all the configurable attributes of a Client. None of
// the options are mandatory, although if you do not want to configure the
// client then you should use the DefaultClient.
type ClientOptions struct {
	// BaseURL is prepended to the URI of all requests made by the Client.
	BaseURL string

	// PreRequestHooks is a list of RequestHooks that get applied in order to
	// requests made by the Client just before they are sent.
	PreRequestHooks []RequestHook

	PostRequestHooks []ResponseHook

	// Transport is the internal HTTP client used to perform the requests.
	Transport *http.Client
}

// New returns a snorlax client configured with the provided ClientOptions.
func NewClient(opts ClientOptions) Client {
	c := client{defaultOptions}

	if opts.BaseURL != "" {
		c.opts.BaseURL = opts.BaseURL
	}

	if opts.PreRequestHooks != nil {
		c.opts.PreRequestHooks = opts.PreRequestHooks
	}

	if opts.PostRequestHooks != nil {
		c.opts.PostRequestHooks = opts.PostRequestHooks
	}

	if opts.Transport != nil {
		c.opts.Transport = opts.Transport
	}

	return &c
}

// Client defines a stateful REST client able to perform HTTP requests.
type client struct {
	opts ClientOptions
}

func (c *client) call(ctx context.Context, method, target string,
	query url.Values, body io.Reader, hooks ...RequestHook) (*Response,
	error) {

	u := strings.Join([]string{c.opts.BaseURL, target}, "")
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
	for _, hook := range append(c.opts.PreRequestHooks, hooks...) {
		hook(req)
	}

	res, err := c.opts.Transport.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform http request: %w", err)
	}

	for _, hook := range c.opts.PostRequestHooks {
		hook(res)
	}

	return &Response{*res}, nil
}

// Delete performs a DELETE request.
func (c *client) Delete(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodDelete, target, query, body, opts...)
}

// Get performs a GET request.
func (c *client) Get(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodGet, target, query, nil, opts...)
}

// Head performs a HEAD request.
func (c *client) Head(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodHead, target, query, nil, opts...)
}

// Options performs a OPTIONS request.
func (c *client) Options(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodOptions, target, query, nil, opts...)
}

// Post performs a POST request.
func (c *client) Post(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodPost, target, query, body, opts...)
}

// Put performs a PUT request.
func (c *client) Put(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodPut, target, query, body, opts...)
}
