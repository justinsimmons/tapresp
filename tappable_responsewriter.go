package tapresp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// A wrapper for http.ResponseWriter that provides access to the underlying
// values in the HTTP Response.
type TappableResponseWriter struct {
	status int
	body   bytes.Buffer

	w http.ResponseWriter
}

// Status code of the HTTP response held by the ResponseWriter.
// Status code will be -1 if it has not yet been set.
func (rw *TappableResponseWriter) StatusCode() int {
	return rw.status
}

// HTTP headers of the response held by the ResponseWriter.
func (rw *TappableResponseWriter) Header() http.Header {
	return rw.w.Header()
}

// Returns a copy of the value held in the response body.
// Modification of the resultant byte slice will in no way
// modify the response body. Instead use Write(...).
func (rw *TappableResponseWriter) Body() []byte {
	return rw.body.Bytes()
}

// Write the HTTP status code to the Reponse
func (rw *TappableResponseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode

	rw.w.WriteHeader(statusCode)
}

// Write data to the response body, note that each write overrides
// previously held in the response body.
func (rw *TappableResponseWriter) Write(data []byte) (int, error) {
	// Clear out anything previously stored in the buffer
	rw.body.Reset()

	n, err := rw.body.Write(data)
	if err != nil {
		return n, err
	}

	if n != len(data) {
		return n, fmt.Errorf("failed to write entirety of input to body")
	}

	return rw.w.Write(data)
}

// Provides a copy of the underlying HTTP response that will be created by
// the http.ResponseWriter.
// Modification of the resultant http.Response struct will in no way
// affect the state of the response writer.
func (rw *TappableResponseWriter) Response() *http.Response {
	return &http.Response{
		StatusCode: rw.status,
		Header:     rw.w.Header().Clone(),
		// Will be a copy of the byte slice held by the buffer
		Body: io.NopCloser(bytes.NewBuffer(rw.body.Bytes())),
	}
}

// Constructs a tappable http.ResponseWriter from another http.ResponseWriter.
func New(w http.ResponseWriter) *TappableResponseWriter {
	var buf bytes.Buffer

	return &TappableResponseWriter{
		w:      w,
		status: 0,
		body:   buf,
	}
}
