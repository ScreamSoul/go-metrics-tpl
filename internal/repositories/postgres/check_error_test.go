package postgres

import (
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIsTemporaryConnectionError(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		inputErr error
		expected bool
	}{
		{
			name:     "Temporary connection error",
			inputErr: &pgconn.PgError{Code: pgerrcode.ConnectionException},
			expected: true,
		},
		{
			name:     "Non-temporary connection error",
			inputErr: &pgconn.PgError{Code: pgerrcode.SyntaxError},
			expected: false,
		},
		{
			name:     "Non-PgError error",
			inputErr: errors.New("some other error"),
			expected: false,
		},
		{
			name:     "Nil error",
			inputErr: nil,
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsTemporaryConnectionError(tc.inputErr)
			assert.Equal(t, tc.expected, result)
		})
	}
}
