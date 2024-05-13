package client

import (
	"errors"
	"net"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func IsTemporaryNetworkError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	var respErr *resty.ResponseError
	if errors.As(err, &respErr) {
		if respErr.Response.StatusCode() == http.StatusRequestTimeout ||
			respErr.Response.StatusCode() == http.StatusServiceUnavailable {
			return true
		}
	}

	return false
}
