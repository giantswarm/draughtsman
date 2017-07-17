package http

import (
	"net/http"
)

// Client is an interface for HTTP clients to implement.
type Client interface {
	Do(*http.Request) (*http.Response, error)
}
