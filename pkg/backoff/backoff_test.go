package backoff

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryWithBackoff(t *testing.T) {
	failingFn := func() error {
		return errors.New("simulated failure")
	}

	successfulFn := func() error {
		return nil
	}

	tests := []struct {
		name             string
		backoffIntervals []time.Duration
		shouldRetry      func(error) bool
		fn               func() error
		expectedError    error
	}{
		{
			name:             "Retry on failure",
			backoffIntervals: []time.Duration{100 * time.Millisecond, 200 * time.Millisecond},
			shouldRetry:      func(err error) bool { return true },
			fn:               failingFn,
			expectedError:    errors.New("simulated failure"),
		},
		{
			name:             "Success on first attempt",
			backoffIntervals: []time.Duration{100 * time.Millisecond, 200 * time.Millisecond},
			shouldRetry:      func(err error) bool { return true },
			fn:               successfulFn,
			expectedError:    nil,
		},
		{
			name:             "No retry on non-retryable error",
			backoffIntervals: []time.Duration{100 * time.Millisecond, 200 * time.Millisecond},
			shouldRetry:      func(err error) bool { return false },
			fn:               failingFn,
			expectedError:    errors.New("simulated failure"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RetryWithBackoff(tt.backoffIntervals, tt.shouldRetry, tt.fn)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
