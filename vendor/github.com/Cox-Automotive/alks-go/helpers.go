package alks

import (
	"net/http"
)

const requestIDHeader = "X-Request-ID"

// GetRequestID returns the ALKS Request ID Header if present or ""
func GetRequestID(resp *http.Response) string {
	return resp.Header.Get(requestIDHeader)
}
