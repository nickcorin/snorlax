package snorlax

import "net/http"

type (
	// RequestHook is a middleware function that can be applied to an HTTP
	// request before it's sent.
	RequestHook func(Client, *http.Request) error
)

// WithBasicAuth sets basic authentication on the request.
func WithBasicAuth(username, password string) RequestHook {
	return func(c Client, r *http.Request) error {
		r.SetBasicAuth(username, password)
		return nil
	}
}

// WithHeader adds a the header key value pair to the request.
func WithHeader(key, value string) RequestHook {
	return func(c Client, r *http.Request) error {
		r.Header.Set(key, value)
		return nil
	}
}
