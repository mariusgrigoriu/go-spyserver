package spyserver

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// SpyServer implements the http.RoundTripper interface to serve a custom
// http.Response or error while capturing the request for later inspection
// A pointer to this type can be assigned to http.Client.Transport to
// override the default RoundTripper
type SpyServer struct {
	Response *http.Response
	Error    error

	request *http.Request
}

// RoundTrip implements the http.RoundTripper interface
func (s *SpyServer) RoundTrip(req *http.Request) (*http.Response, error) {
	s.request = req
	return s.Response, s.Error
}

// GetRequest returns the http.Request that was passed during the last RoundTrip
func (s *SpyServer) GetRequest() *http.Request {
	return s.request
}

// CloseDetector is an io.ReadCloser that tracks whether Close has been called.
// This struct is useful because the behavior of Close after the first call is undefined.
type CloseDetector struct {
	io.ReadCloser
	closed bool
}

// NewCloseDetector wraps an existing io.ReadCloser.
func NewCloseDetector(rc io.ReadCloser) *CloseDetector {
	var closed bool
	return &CloseDetector{rc, closed}
}

// NewCloseDetectorFromString wraps a NopCloser containing the string
func NewCloseDetectorFromString(s string) *CloseDetector {
	return NewCloseDetector(ioutil.NopCloser(strings.NewReader(s)))
}

// Read is a pass-through function that delegates to the wrapped Read function.
func (cd *CloseDetector) Read(p []byte) (n int, err error) {
	return cd.ReadCloser.Read(p)
}

// Close delegates to the wrapped Close function and tracks whether the function was called.
func (cd *CloseDetector) Close() error {
	cd.closed = true
	if cd.ReadCloser == nil {
		return nil
	}
	return cd.ReadCloser.Close()
}

// IsClosed returns whether Close was called
func (cd *CloseDetector) IsClosed() bool {
	return cd.closed
}
