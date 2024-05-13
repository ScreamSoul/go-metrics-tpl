//go:build ignore

package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func GetPublicFields(s interface{}) (<-chan reflect.StructField, <-chan interface{}, error) {
	fields := make(chan reflect.StructField)
	values := make(chan interface{})

	if s == nil {
		return fields, values, errors.New("s is nil")
	}

	val := reflect.Indirect(reflect.ValueOf(s))
	if val.Type().Kind() != reflect.Struct {
		return fields, values, errors.New("s is not a struct")
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered:", r)
			}
		}()
		defer close(fields)
		defer close(values)

		typeOfS := val.Type()

		numField := val.NumField()
		for i := 0; i < numField; i++ {
			f := val.Field(i)

			if !f.CanInterface() {
				continue
			}

			field := typeOfS.Field(i)

			if field.PkgPath == "" {
				fields <- field
				values <- f.Interface()
			}
		}
	}()

	return fields, values, nil
}
