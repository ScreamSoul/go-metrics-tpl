package client

import (
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestIsTemporaryNetworkError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "net.Error timeout",
			err: &net.OpError{
				Op:     "dial",
				Net:    "tcp",
				Source: nil,
				Addr:   nil,
				Err:    &net.DNSError{IsTimeout: true}, // Correctly simulating a timeout error
			},
			want: true,
		},
		{
			name: "resty.ResponseError request timeout",
			err: &resty.ResponseError{
				Response: &resty.Response{
					Request: &resty.Request{},
					RawResponse: &http.Response{
						StatusCode: http.StatusRequestTimeout,
					},
				},
			},
			want: true,
		},
		{
			name: "resty.ResponseError service unavailable",
			err: &resty.ResponseError{
				Response: &resty.Response{
					Request: &resty.Request{},
					RawResponse: &http.Response{
						StatusCode: http.StatusServiceUnavailable,
					},
				},
			},
			want: true,
		},
		{
			name: "resty.ResponseError other status code",
			err: &resty.ResponseError{
				Response: &resty.Response{
					Request: &resty.Request{},
					RawResponse: &http.Response{
						StatusCode: http.StatusOK,
					},
				},
			},
			want: false,
		},
		{
			name: "other error",
			err:  errors.New("some other error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTemporaryNetworkError(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
