package transit

import "net/http"

var defaultOptions = clientOptions{
	adapters: make([]Adapter, 0),
	headers: make(http.Header),
	transport: http.DefaultClient,
}

type clientOptions struct {
	adapters  []Adapter
	baseURL   string
	headers   http.Header
	transport *http.Client
}

// Adapter defines a functional wrapper type to preprocess the http request
// before it gets sent.
type Adapter interface {
	Adapt(*http.Request)
}

// ClientOption defines a functional wrapper type to provide configuration
// options for the transit Client.
type ClientOption func(*clientOptions)

// WithAdapter provides a ClientOption to configure adapters which will be
// sequentially run before executing requests.
func WithAdapter(adapter Adapter) ClientOption {
	return func(opts *clientOptions) {
		opts.adapters = append(opts.adapters, adapter)
	}
}

// WithBaseURL provides a ClientOption to configure a global base URL for the
// Client. If a base URL is set, it will prefix all request paths made by the
// client.
func WithBaseURL(url string) ClientOption {
	return func(opts *clientOptions) {
		opts.baseURL = url
	}
}

// WithHeader provides a ClientOption to configure request headers to be
// included with each request made by a transit client.
func WithHeader(key, value string) ClientOption {
	return func(opts *clientOptions) {
		opts.headers.Set(key, value)
	}
}

// WithTransport provides a ClientOption to configure the internal http.Client
// transport of the transit Client.
func WithTransport(t *http.Client) ClientOption {
	return func(opts *clientOptions) {
		opts.transport = t
	}
}
