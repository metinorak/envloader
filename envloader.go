package envloader

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type EnvLoader interface {
	Load(model any) error
}

type envLoader struct {
	envReader EnvReader
}

// Creates a new instance of EnvLoader
// Example:
//
//	loader := envloader.New()
func New() EnvLoader {
	return &envLoader{
		envReader: &DefaultEnvReader{},
	}
}

func (el *envLoader) toSnakeUpperCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}

func (el *envLoader) loadFromEnvToMap(envValue string, fieldValue reflect.Value) error {
	pairs := strings.Split(envValue, ",")

	mapValue := reflect.MakeMap(fieldValue.Type())

	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			return fmt.Errorf("invalid map value: %s", envValue)
		}

		// set the key and value regarding the value type
		switch mapValue.Type().Elem().Kind() {
		case reflect.String:
			mapValue.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(kv[1]))
		case reflect.Int:
			intValue, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			mapValue.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(intValue))
		case reflect.Float64:
			floatValue, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				return err
			}
			mapValue.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(floatValue))
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(kv[1])
			if err != nil {
				return err
			}
			mapValue.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(boolValue))
		default:
			return fmt.Errorf("unsupported map value type: %s", mapValue.Type().Elem().Kind())
		}
	}

	fieldValue.Set(mapValue)

	return nil
}

func (el *envLoader) loadFromEnvToModel(keyPrefix string, model any) error {
	value := reflect.ValueOf(model).Elem()
	valueType := value.Type()

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)

		var key string

		// if the field has a tag, use it
		if tag, ok := field.Tag.Lookup("env"); ok {
			key = tag
		} else {
			// otherwise, use the field name
			key = el.toSnakeUpperCase(field.Name)
		}

		kindOfValue := value.Field(i).Kind()
		fieldValue := value.Field(i)

		var currentKey string
		if keyPrefix == "" {
			currentKey = key
		} else {
			currentKey = fmt.Sprintf("%s%s%s", keyPrefix, "_", key)
		}

		envValue, envExists := el.envReader.LookupEnv(currentKey)

		if !envExists && kindOfValue != reflect.Struct {
			// if default tag exists, use it
			if defaultTag, ok := field.Tag.Lookup("default"); ok {
				envValue = defaultTag
			}
		}

		switch kindOfValue {
		case reflect.String:
			fieldValue.SetString(envValue)

		case reflect.Int:
			intValue, err := strconv.Atoi(envValue)
			if err != nil {
				return err
			}
			fieldValue.SetInt(int64(intValue))

		case reflect.Float64:
			floatValue, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return err
			}
			fieldValue.SetFloat(floatValue)

		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envValue)
			if err != nil {
				return err
			}
			fieldValue.SetBool(boolValue)

		case reflect.Slice:
			sliceValue := strings.Split(envValue, ",")
			fieldValue.Set(reflect.ValueOf(sliceValue))

		case reflect.Map:
			err := el.loadFromEnvToMap(envValue, fieldValue)
			if err != nil {
				return err
			}

		case reflect.Struct:
			el.loadFromEnvToModel(currentKey, fieldValue.Addr().Interface())
		}
	}

	return nil

}

func (el *envLoader) loadFromEnv(model any) error {
	// check the model type
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("model must be a pointer")
	}

	if reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("model must be a pointer to a struct")
	}

	// find all env keys and set to model
	el.loadFromEnvToModel("", model)

	return nil
}

// Loads the environment variables into the provided model
func (el *envLoader) Load(model any) error {
	err := el.loadFromEnv(model)
	if err != nil {
		return err
	}

	return nil
}
