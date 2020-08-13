package snorlax

import "net/http"

// ClientOption defines an option to configure all requests made by the client.
type ClientOption interface {
	Apply(*Client)
}

// ClientOptionFunc is a function to ClientOptionFunc adapter.
type ClientOptionFunc func(*Client)

// Apply satisfies the ClientOption interface.
func (f ClientOptionFunc) Apply(c *Client) {
	f(c)
}

// WithBaseURL returns a ClientOptionFunc to configure the base url of the
// client.
func WithBaseURL(url string) ClientOptionFunc {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithRequestOptions returns a ClientOptionFunc to set RequestOptions to be
// applied to all requests.
func WithRequestOptions(opts ...RequestOption) ClientOptionFunc {
	return func(c *Client) {
		c.requestOptions = append(c.requestOptions, opts...)
	}
}

// WithTransport provides a RequestOptionFunc to configure the internal
// http.Client transport of the transit Client.
func WithTransport(t *http.Client) ClientOptionFunc {
	return func(c *Client) {
		c.transport = t
	}
}

// RequestOption defines an option to configure individual requests made by the
// client. You can set RequestOptions that you want to be applied to all
// requests made by the client by using WithRequestOptions.
type RequestOption interface {
	Apply(*http.Request)
}

// RequestOptionFunc is a function to RequestOptionFunc adapter.
type RequestOptionFunc func(*http.Request)

// Apply satisfies the RequestOption interface.
func (f RequestOptionFunc) Apply(r *http.Request) {
	f(r)
}

// WithBasicAuth returns a RequestOptionFunc to set basic authentication on a
// request.
func WithBasicAuth(username, password string) RequestOptionFunc {
	return func(r *http.Request) {
		r.SetBasicAuth(username, password)
	}
}

// WithHeader provides a RequestOptionFunc to configure request headers to be
// included with each request made by a transit client.
func WithHeader(key, value string) RequestOptionFunc {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}
