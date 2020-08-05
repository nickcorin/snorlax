package transit

import (
	"io/ioutil"
	"net/http"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// Copy all request headers to the response.
	for k, v := range r.Header {
		for _, h := range v {
			w.Header().Add(k, h)
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(requestBody)
}
