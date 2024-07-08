package logging

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggingSuite struct {
	suite.Suite
	logger *zap.Logger
}

func (s *LoggingSuite) SetupTest() {
	// Reset the logger to the default NOP logger before each test
	err := Initialize("info")
	s.Require().NoError(err)
	s.logger = GetLogger()
}

func (s *LoggingSuite) TestInitialize() {
	testCases := []struct {
		name     string
		level    string
		expected zapcore.Level
	}{
		{"Debug Level", "debug", zapcore.DebugLevel},
		{"Info Level", "info", zapcore.InfoLevel},
		{"Warn Level", "warn", zapcore.WarnLevel},
		{"Error Level", "error", zapcore.ErrorLevel},
		{"DPanic Level", "dpanic", zapcore.DPanicLevel},
		{"Panic Level", "panic", zapcore.PanicLevel},
		{"Fatal Level", "fatal", zapcore.FatalLevel},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := Initialize(tc.level)
			s.Require().NoError(err, "Initialize should not return an error")
			s.True(GetLogger().Core().Enabled(tc.expected), "Logger level should be enabled for the expected level")
		})
	}
}

func (s *LoggingSuite) TestInitialize_Fail() {
	s.Run("Error parse level", func() {
		err := Initialize("fali lavel")
		s.Require().Error(err)
	})

}

func TestLoggingSuite(t *testing.T) {
	suite.Run(t, new(LoggingSuite))
}
