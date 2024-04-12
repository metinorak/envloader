package goenv

import (
	"reflect"
	"strings"
)

type structField reflect.StructField

func (sf structField) isRequired() bool {
	if tag, ok := sf.Tag.Lookup("required"); ok && tag == "true" {
		return true
	}

	return false
}

func (sf structField) getDefaultValue() (string, bool) {
	return sf.Tag.Lookup("default")
}

func (sf structField) toSnakeUpperCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}

func (sf structField) getEnvName() string {
	var key string

	if tag, ok := sf.Tag.Lookup("env"); ok {
		key = tag
	} else {
		// otherwise, use the field name
		key = sf.toSnakeUpperCase(sf.Name)
	}

	return key
}
