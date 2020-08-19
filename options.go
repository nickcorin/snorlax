package snorlax

import "net/http"

var defaultOptions = ClientOptions{
	CallOptions: make([]CallOption, 0),
	Transport:   http.DefaultClient,
}

type ClientOptions struct {
	BaseURL     string
	CallOptions []CallOption
	Transport   *http.Client
}

// CallOption defines an option to configure individual requests made by the
// client. You can set CallOptions that you want to be applied to all
// requests made by the client by adding them to the Client's CallOptions.
type CallOption interface {
	Apply(*http.Request)
}

// CallOptionFunc is a function to CallOptionFunc adapter.
type CallOptionFunc func(*http.Request)

// Apply satisfies the CallOption interface.
func (f CallOptionFunc) Apply(r *http.Request) {
	f(r)
}

// WithBasicAuth returns a CallOptionFunc to set basic authentication on a
// request.
func WithBasicAuth(username, password string) CallOptionFunc {
	return func(r *http.Request) {
		r.SetBasicAuth(username, password)
	}
}

// WithHeader provides a CallOptionFunc to set a header on a request.
func WithHeader(key, value string) CallOptionFunc {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}
