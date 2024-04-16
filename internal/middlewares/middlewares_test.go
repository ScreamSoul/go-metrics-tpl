package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipDecompressMiddleware(t *testing.T) {
	middleware := GzipDecompressMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body) // Use io.ReadAll instead of ioutil.ReadAll
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	}))

	tests := []struct {
		name                 string
		contentEncoding      string
		expectedResponseBody string
		encodingFunc         func(body *bytes.Buffer, data string)
	}{
		{
			name:                 "Gzip Encoded",
			contentEncoding:      "gzip",
			expectedResponseBody: "Hello, World!",
			encodingFunc: func(body *bytes.Buffer, data string) {
				gw := gzip.NewWriter(body)
				gw.Write([]byte(data))
				gw.Close()
			},
		},
		{
			name:                 "Not Gzip Encoded",
			contentEncoding:      "",
			expectedResponseBody: "Hello, World!",
			encodingFunc: func(body *bytes.Buffer, data string) {
				body.Write([]byte(data))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer

			tt.encodingFunc(&body, tt.expectedResponseBody)

			req, err := http.NewRequest("POST", "/", &body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Encoding", tt.contentEncoding)

			// Record the response
			rr := httptest.NewRecorder()
			middleware.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, http.StatusOK, rr.Code)

			// Check the response body
			assert.Equal(t, tt.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestGzipCompressMiddleware(t *testing.T) {
	middleware := GzipCompressMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	tests := []struct {
		name                 string
		acceptEncodingHeader string
		contentTypeHeader    string
		expectGzip           bool
	}{
		{
			name:                 "Gzip Supported and Content Type Supported",
			acceptEncodingHeader: "gzip",
			contentTypeHeader:    "application/json",
			expectGzip:           true,
		},
		{
			name:                 "Gzip Supported but Content Type Not Supported",
			acceptEncodingHeader: "gzip",
			contentTypeHeader:    "image/png",
			expectGzip:           false,
		},
		{
			name:                 "Gzip Not Supported",
			acceptEncodingHeader: "",
			contentTypeHeader:    "application/json",
			expectGzip:           false,
		},
		{
			name:                 "Gzip Supported and Content Type Not Specified",
			acceptEncodingHeader: "gzip",
			contentTypeHeader:    "",
			expectGzip:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", "/", nil)
			assert.NoError(t, err)

			// Set headers
			req.Header.Set("Accept-Encoding", tt.acceptEncodingHeader)
			req.Header.Set("Content-Type", tt.contentTypeHeader)

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Serve the request
			middleware.ServeHTTP(rr, req)

			// Check if the response is gzipped
			if tt.expectGzip {
				assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

				// Decompress the response body
				gr, err := gzip.NewReader(rr.Body)
				assert.NoError(t, err)
				defer gr.Close()

				body, err := io.ReadAll(gr)
				assert.NoError(t, err)
				assert.Equal(t, "Hello, World!", string(body))
			} else {
				assert.Empty(t, rr.Header().Get("Content-Encoding"))
				assert.Equal(t, "Hello, World!", rr.Body.String())
			}
		})
	}
}
