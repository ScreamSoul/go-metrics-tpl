package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitialize(t *testing.T) {
	originalLogger := GetLogger()
	defer func() { log = originalLogger }()

	log = nil

	// некорректный уровень логирования
	err := Initialize("invalid_level")
	assert.NotNil(t, err, "Expected an error when initializing with an invalid level")
	assert.Nil(t, log, "Expected logger to not be initialized")

	// корректный уровень логирования
	err = Initialize("info")
	assert.Nil(t, err, "Expected no error when initializing with a valid level")
	assert.NotNil(t, log, "Expected logger to be initialized")
	assert.False(t, log.Core().Enabled(zap.DebugLevel), "The debug level should not be supported")
	assert.True(t, log.Core().Enabled(zap.InfoLevel), "The info level should be supported")
}
