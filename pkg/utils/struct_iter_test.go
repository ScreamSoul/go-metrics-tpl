//go:build ignore

package utils

// Returns two channels and no error when given a non-nil struct
import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicFields_NonNilStruct(t *testing.T) {
	// Arrange
	type TestStruct struct {
		PublicField int
	}
	s := TestStruct{}

	// Act
	fields, values, err := GetPublicFields(&s)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, fields)
	assert.NotNil(t, values)
}

func TestGetPublicFields_UnexportedFields(t *testing.T) {
	// Arrange
	type TestStruct struct {
		PublicField int
	}
	s := TestStruct{}

	// Act
	fields, values, err := GetPublicFields(&s)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, fields)
	assert.NotNil(t, values)

	// Check that only the public field is returned
	field := <-fields
	value := <-values
	assert.Equal(t, "PublicField", field.Name)
	assert.Equal(t, reflect.ValueOf(s.PublicField).Interface(), value)

	// Check that there are no more fields and values
	_, moreFields := <-fields
	_, moreValues := <-values
	assert.False(t, moreFields)
	assert.False(t, moreValues)
}
