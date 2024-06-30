package server_test

import (
	"os"
	"testing"
	"time"

	"github.com/screamsoul/go-metrics-tpl/internal/server"
	"github.com/stretchr/testify/assert"
)

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

			cfg, err := server.NewConfig()
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBackoffIntervals, cfg.Postgres.BackoffIntervals)
			assert.Equal(t, tt.expectedBackoffRetries, cfg.Postgres.BackoffRetries)

			for k := range tt.envVars {
				t.Setenv(k, "")
			}
		})
	}
}
