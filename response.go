package snorlax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Response is a type alias for http.Response.
type Response struct {
	http.Response
}

// IsSuccess returns whether the response code is within the 2XX range.
func (r *Response) IsSuccess() bool {
	return r.StatusCode < http.StatusMultipleChoices
}

// JSON reads and unmarshals the response body into out.
func (r *Response) JSON(out interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err = json.Unmarshal(body, &out); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}

// RawBody returns an io.Reader containing the data returned in the response
// body.
func (r *Response) RawBody() (io.Reader, error) {
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return bytes.NewBuffer(data), nil
}
