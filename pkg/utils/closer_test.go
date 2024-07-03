// Closer successfully closes without errors
package utils_test

import (
	"testing"

	"github.com/screamsoul/go-metrics-tpl/pkg/utils"
	"github.com/stretchr/testify/mock"
)

type MockCloser struct {
	mock.Mock
}

func (m *MockCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestCloseForse_SuccessfulClose(t *testing.T) {
	mockCloser := new(MockCloser)
	mockCloser.On("Close").Return(nil)

	utils.CloseForse(mockCloser)

	mockCloser.AssertExpectations(t)
}
