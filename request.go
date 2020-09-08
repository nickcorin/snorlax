package snorlax

import "net/http"

type (
	// RequestHook is a middleware function that can be applied to an HTTP
	// request before it's sent.
	RequestHook func(*http.Request)

	// ResponseHook is a middleware function that can be applied to an HTTP
	// response as it's received.
	ResponseHook func(*http.Response)
)

// WithBasicAuth returns a RequestHook to set basic authentication on a
// request.
func WithBasicAuth(username, password string) RequestHook {
	return func(r *http.Request) {
		r.SetBasicAuth(username, password)
	}
}

// WithHeader provides a RequestHook to set a header on a request.
func WithHeader(key, value string) RequestHook {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}
