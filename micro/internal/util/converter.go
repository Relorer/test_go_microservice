package util

import (
	"reflect"
	"unicode"
)

// Recursively converts an object into a map, where the keys are the names of the fields
func ToMaps(obj interface{}) interface{} {
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Struct:
		typ := value.Type()
		result := make(map[string]interface{}, typ.NumField())

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			jsonKey := field.Tag.Get("json")

			// privacy check
			if unicode.IsLower([]rune(field.Name)[0]) {
				continue
			}
			if jsonKey == "" {
				jsonKey = field.Name
			}
			result[jsonKey] = ToMaps(value.Field(i).Interface())
		}
		return result

	case reflect.Slice, reflect.Array:
		result := make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			result[i] = ToMaps(value.Index(i).Interface())
		}
		return result

	default:
		return value.Interface()
	}
}
