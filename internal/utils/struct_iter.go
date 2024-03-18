package utils

import "reflect"

func isPublicField(field reflect.StructField) bool {
	// Проверяем, является ли имя поля публичным (начинается с заглавной буквы)
	return field.Name[0] >= 'A' && field.Name[0] <= 'Z'
}

func PublicFields(s interface{}) (<-chan reflect.StructField, <-chan interface{}) {
	fields := make(chan reflect.StructField)
	values := make(chan interface{})

	go func() {
		defer close(fields)
		defer close(values)

		val := reflect.ValueOf(s)

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		typeOfS := val.Type()

		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)

			if f.Kind() == reflect.Struct {
				continue
			}

			field := typeOfS.Field(i)

			if isPublicField(field) {
				fields <- field
				values <- f.Interface()
			}
		}
	}()

	return fields, values
}
