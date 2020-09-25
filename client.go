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

	"github.com/sirupsen/logrus"
)

// Client defines a wrapper around an http.Client making it easier to send
// requests to RESTful APIs.
type Client interface {
	// AddHeader appends a header value to the client to be sent in every
	// request. To replace the current existing header use SetHeader.
	AddHeader(key, value string) Client

	// AddRequestHook appends a RequestHook to the list of hooks which are to be
	// run just before the client sends a request. RequestHooks are executed in
	// the order they are added.
	AddRequestHook(hook RequestHook) Client

	// AddRequestHooks is a convenience function which calls AddRequestHook
	// multiple times.
	AddRequestHooks(hooks ...RequestHook) Client

	// Get performs a Get request. You can optionally configure the request
	// using RequestHooks, or by configuring the client if you need to configure
	// all requests.
	Delete(ctx context.Context, target string, query url.Values, body io.Reader,
		hooks ...RequestHook) (*Response, error)

	// Get performs a Get request. You can optionally configure the request
	// using RequestHooks, or by configuring the client if you need to configure
	// all requests.
	Get(ctx context.Context, target string, query url.Values,
		hooks ...RequestHook) (*Response, error)

	// Head performs a Head request. You can optionally configure the request
	// using RequestHooks, or by configuring the client if you need to configure
	// all requests.
	Head(ctx context.Context, target string, query url.Values,
		hooks ...RequestHook) (*Response, error)

	// Logger returns the client's internal logger.
	Logger() *logrus.Logger

	// Options performs a Options request. You can optionally configure the
	// request using RequestHooks, or by configuring the client if you need to
	// configure  all requests.
	Options(ctx context.Context, target string, query url.Values,
		hooks ...RequestHook) (*Response, error)

	// Post performs a Post request. You can optionally configure the request
	// using RequestHooks, or by configuring the client if you need to configure
	// all requests.
	Post(ctx context.Context, target string, query url.Values, body io.Reader,
		hooks ...RequestHook) (*Response, error)

	// Put performs a Put request. You can optionally configure the request
	// using RequestHooks, or by configuring the client if you need to configure
	// all requests.
	Put(ctx context.Context, target string, query url.Values, body io.Reader,
		hooks ...RequestHook) (*Response, error)

	// RemoveProxy removes any currently set proxy URL in the Client's
	// transport.
	RemoveProxy() Client

	// SetBaseURL sets a host URL inside the client which is prepended to all
	// request URLs performed by the Client.
	SetBaseURL(url string) Client

	// SetHeader sets a header value in the client to be sent in every request.
	// This will overwrite any exiting headers present associated with the same
	// key. To add headers to the key instead of replacing them use AddHeader.
	SetHeader(key, value string) Client

	// SetHTTPClient replaces the internal http.Client that Snorlax uses to
	// perform requests. Use this if you need finer control over the client's
	// internals.
	SetHTTPClient(c *http.Client) Client

	// SetLogLevel sets the amount of logs the client will produce. The lower
	// the level, the less logs will be written. By default, Snorlax uses the
	// lowest possible level - PanicLevel.
	SetLogLevel(level logrus.Level) Client

	// SetProxy sets the proxy URL in the clent's transport. If the URL fails
	// to parse, nothing is set. This function fails silently. If you need more
	// of a guarantee rather create your own http.Client with your proxy set and
	// use SetHTTPClient.
	SetProxy(url string) Client
}

// DefaultClient is a Snorlax client configured with all of the default options.
var DefaultClient = &client{
	opts: Defaults(),
}

// Delete performs a delete request using the DefaultClient. You can optionally
// configure the request using RequestHooks. If you need to configure every
// request then consider not using the DefaultClient.
func Delete(ctx context.Context, target string, query url.Values,
	body io.Reader, hooks ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodDelete, target, query, body,
		hooks...)
}

// Get performs a get request using the DefaultClient. You can optionally
// configure the request using RequestHooks. If you need to configure every
// request then consider not using the DefaultClient.
func Get(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodGet, target, query, nil, opts...)
}

// Head performs a head request using the DefaultClient. You can optionally
// configure the request using RequestHooks. If you need to configure every
// request then consider not using the DefaultClient.
func Head(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodHead, target, query, nil, opts...)
}

// Options performs an options request using the DefaultClient. You can
// optionally configure the request using RequestHooks. If you need to configure
// every request then consider not using the DefaultClient.
func Options(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodOptions, target, query, nil,
		opts...)
}

// Post performs a post request using the DefaultClient. You can optionally
// configure the request using RequestHooks. If you need to configure every
// request then consider not using the DefaultClient.
func Post(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodPost, target, query, body,
		opts...)
}

// Post performs a post request using the DefaultClient. You can optionally
// configure the request using RequestHooks. If you need to configure every
// request then consider not using the DefaultClient.
func Put(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return DefaultClient.call(ctx, http.MethodPut, target, query, body, opts...)
}

// NewClient constructs a new Client configured with the provided ClientOptions.
func NewClient(opts *ClientOptions) Client {
	return &client{opts}
}

type client struct {
	opts *ClientOptions
}

// ClientOptions contains the configuration options for a Snorlax client.
type ClientOptions struct {
	BaseURL     string
	WithMetrics bool

	headers      http.Header
	httpClient   *http.Client
	logger       *logrus.Logger
	logLevel     logrus.Level
	proxyURL     *url.URL
	requestHooks []RequestHook
}

// Defaults returns a set of default ClientOptions.
func Defaults() *ClientOptions {
	opts := ClientOptions{
		BaseURL:     "",
		WithMetrics: false,

		headers:      make(http.Header),
		httpClient:   http.DefaultClient,
		logger:       logrus.New(),
		logLevel:     logrus.PanicLevel,
		proxyURL:     nil,
		requestHooks: make([]RequestHook, 0),
	}

	opts.logger.SetLevel(opts.logLevel)
	opts.logger.SetFormatter(&logrus.JSONFormatter{})

	return &opts
}

func (c *client) call(ctx context.Context, method, target string,
	query url.Values, body io.Reader, hooks ...RequestHook) (*Response,
	error) {

	u := strings.Join([]string{c.opts.BaseURL, target}, "")
	uri, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %s: %w", uri, err)
	}

	if uri.RawQuery != "" {
		c.opts.logger.Warnln("query parameters should not be set on the " +
			"path, they will be overridden")
	}
	uri.RawQuery = query.Encode()

	c.opts.logger.Tracef("uri parsed as %s", uri.String())

	req, err := http.NewRequestWithContext(ctx, method, uri.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	// Set the request headers with all the headers configured in the client.
	if c.opts.headers == nil {
		c.opts.headers = make(http.Header)
	}
	req.Header = c.opts.headers

	c.opts.logger.Tracef("setting content-length header to %d",
		req.ContentLength)
	// Automatically add the Content-Length header.
	req.Header.Set(http.CanonicalHeaderKey("Content-Length"),
		strconv.FormatInt(req.ContentLength, 10))

	// httpClient is usually nil on the first request made by the client. This
	// prevents panics by using the http.DefaultClient. In most cases, this will
	// be sufficient. In cases where the caller wants more control over the
	// client's configuration - SethttpClient can be used.
	if c.opts.httpClient == nil {
		c.opts.logger.Trace("internal client set to http.DefaultClient")
		c.opts.httpClient = http.DefaultClient
	}

	c.opts.logger.Trace("running pre-request hooks")
	// We first apply the request options from the client, so that they can be
	// optionally overridden by individual request options.
	for _, hook := range append(c.opts.requestHooks, hooks...) {
		if err = hook(c, req); err != nil {
			return nil, fmt.Errorf("failed to execute pre-request hook: %w",
				err)
		}
	}
	c.opts.logger.Trace("pre-request hooks complete")

	c.opts.logger.WithField("url", req.URL.String()).Trace("performing request")
	reqStart := time.Now()
	res, err := c.opts.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform http request: %w", err)
	}

	c.opts.logger.WithFields(logrus.Fields{
		"method":      method,
		"latency":     time.Since(reqStart).Seconds(),
		"status_code": res.StatusCode,
		"url":         uri.String(),
	}).Debug("request complete")

	// TODO: Rethink how to provide the target path as a label effectively.
	// clients sending requests to dynamic paths can overload prometheus.
	if c.opts.WithMetrics {
		latencyHist.WithLabelValues(method,
			strconv.Itoa(res.StatusCode), req.URL.Path).Observe(
			time.Since(reqStart).Seconds())
	}

	return &Response{*res}, nil
}

// AddHeader appends a header value to the client to be sent in every request.
// To replace the current existing header use SetHeader.
func (c *client) AddHeader(key, value string) Client {
	c.opts.logger.WithFields(logrus.Fields{"key": key, "value": value}).
		Trace("setting request header")
	c.opts.headers.Add(key, value)
	return c
}

// AddRequestHook appends a RequestHook to the list of hooks which are to be run
// just before the client sends a request. RequestHooks are executed in the
// order they are added.
func (c *client) AddRequestHook(hook RequestHook) Client {
	if c.opts.requestHooks == nil {
		c.opts.requestHooks = make([]RequestHook, 0)
	}

	c.opts.requestHooks = append(c.opts.requestHooks, hook)
	return c
}

// AddRequestHooks is a convenience function which calls AddRequestHook multiple
// times.
func (c *client) AddRequestHooks(hooks ...RequestHook) Client {
	for _, hook := range hooks {
		c.AddRequestHook(hook)
	}

	return c
}

// Delete satisfies the Client interface.
func (c *client) Delete(ctx context.Context, target string, query url.Values,
	body io.Reader, hooks ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodDelete, target, query, body, hooks...)
}

// Get satisfies the Client interface.
func (c *client) Get(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodGet, target, query, nil, opts...)
}

// Head satisfies the Client interface.
func (c *client) Head(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodHead, target, query, nil, opts...)
}

// Logger satisfies the Client interface.
func (c *client) Logger() *logrus.Logger {
	return c.opts.logger
}

// Options satisfies the Client interface.
func (c *client) Options(ctx context.Context, target string, query url.Values,
	opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodOptions, target, query, nil, opts...)
}

// Post satisfies the Client interface.
func (c *client) Post(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodPost, target, query, body, opts...)
}

// Put satisfies the Client interface.
func (c *client) Put(ctx context.Context, target string, query url.Values,
	body io.Reader, opts ...RequestHook) (*Response, error) {
	return c.call(ctx, http.MethodPut, target, query, body, opts...)
}

// RemoveProxy removes the currently set proxy.
func (c *client) RemoveProxy() Client {
	t, ok := c.opts.httpClient.Transport.(*http.Transport)
	if !ok {
		c.opts.logger.Warn("proxy not removed: client transport failed " +
			"assertion")
		return c
	}

	c.opts.proxyURL = nil
	t.Proxy = nil
	c.opts.logger.Trace("proxy removed")

	return c
}

// SetBaseURL sets the url that is prepended to all request URLs.
func (c *client) SetBaseURL(u string) Client {
	if _, err := url.Parse(u); err != nil {
		c.opts.logger.WithField("url", u).
			Warn("base url not set: failed to parse")
		return c
	}

	c.opts.BaseURL = u
	c.opts.logger.WithField("url", u).Trace("base url set")
	return c
}

// SetHeader sets a header value in the client to be sent in every request. This
// will overwrite any exiting headers present associated with the same key. To
// add headers to the key instead of replacing them use AddHeader.
func (c *client) SetHeader(key, value string) Client {
	c.opts.headers.Set(key, value)
	c.opts.logger.WithField("key", key).WithField("value", value).
		Trace("header set")
	return c
}

// SetHTTPClient sets the internal http.client that Snorlax uses to perform
// requests. Use this if you want to configure client internals like timeouts.
func (c *client) SetHTTPClient(client *http.Client) Client {
	c.opts.httpClient = client
	c.opts.logger.Trace("http client set")
	return c
}

// SetLogLevel satisfies the Client interface.
func (c *client) SetLogLevel(level logrus.Level) Client {
	c.opts.logger.SetLevel(level)
	c.opts.logLevel = level
	c.opts.logger.WithField("level", level.String()).Trace("log level set")
	return c
}

// SetRequestHooks sets the RequestHooks to be run just before the client
// performs requests. These are run in order. Calling SetRequestHooks will
// replace any existing RequestHooks that have been added. To add RequestHooks
// without replacing other hooks use AddRequestHook(s).
func (c *client) SetRequestHooks(hooks []RequestHook) Client {
	c.opts.requestHooks = hooks
	return c
}

// SetProxy sets the proxy URL for the Snorlax client. If the provided URL fails
// to be parsed then nothing will be set.
func (c *client) SetProxy(u string) Client {
	t, ok := c.opts.httpClient.Transport.(*http.Transport)
	if !ok {
		c.opts.logger.Warn("proxy not set: client transport failed " +
			"assertion")
		return c
	}

	proxyURL, err := url.Parse(u)
	if err != nil {
		c.opts.logger.WithField("url", u).
			Warn("proxy url not set: failed to parse")
		return c
	}

	c.opts.proxyURL = proxyURL
	t.Proxy = http.ProxyURL(proxyURL)
	c.opts.logger.WithField("url", u).Trace("proxy url set")

	return c
}
