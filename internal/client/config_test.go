package client_test

import (
	"os"
	"testing"
	"time"

	"github.com/screamsoul/go-metrics-tpl/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		config         client.Config
		expectedServer string
		expectedUpdate string
	}{
		{
			name: "Default Config",
			config: client.Config{
				Server: client.Server{
					ListenServerHost: "localhost:8080",
					CompressRequest:  true,
					BackoffIntervals: []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second},
					BackoffRetries:   true,
				},
				ReportInterval: 10,
				PollInterval:   2,
				LogLevel:       "INFO",
			},
			expectedServer: "http://localhost:8080",
			expectedUpdate: "http://localhost:8080/updates/",
		},
		// Add more test cases as needed
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedServer, tc.config.GetServerURL())
			assert.Equal(t, tc.expectedUpdate, tc.config.GetUpdateMetricURL())
		})
	}
}

func TestBackoffIntervalConfig(t *testing.T) {
	testTable := []struct {
		name                     string
		envVars                  map[string]string
		expectedBackoffIntervals []time.Duration
		expectedBackoffRetries   bool
	}{
		{
			name:                     "Default configuration",
			envVars:                  map[string]string{},
			expectedBackoffIntervals: []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second},
			expectedBackoffRetries:   true,
		},
		{
			name: "Custom backoff intervals",
			envVars: map[string]string{
				"BACKOFF_INTERVALS": "1s,2s,3s",
			},
			expectedBackoffIntervals: []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second},
			expectedBackoffRetries:   true,
		},
		{
			name: "Disable backoff retries",
			envVars: map[string]string{
				"BACKOFF_RETRIES": "false",
			},
			expectedBackoffIntervals: nil,
			expectedBackoffRetries:   false,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = nil
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg, err := client.NewConfig()
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBackoffIntervals, cfg.Server.BackoffIntervals)
			assert.Equal(t, tt.expectedBackoffRetries, cfg.Server.BackoffRetries)

			for k := range tt.envVars {
				t.Setenv(k, "")
			}
		})
	}
}
