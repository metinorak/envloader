package goenv

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type ErrParseEnvValue struct {
	Key   string
	Value string
}

func (e *ErrParseEnvValue) Error() string {
	return fmt.Sprintf("failed to parse environment variable %s: %s", e.Key, e.Value)
}

var envReader EnvReader = &DefaultEnvReader{}

func loadFromEnvToMap(envValue string, fieldValue reflect.Value) error {
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

func loadFromEnvToModel(keyPrefix string, model any) error {
	value := reflect.ValueOf(model).Elem()
	valueType := value.Type()

	for i := 0; i < valueType.NumField(); i++ {
		field := structField(valueType.Field(i))
		key := field.getEnvName()

		kindOfValue := value.Field(i).Kind()
		fieldValue := value.Field(i)

		if key == "-" && kindOfValue != reflect.Struct {
			continue
		}

		var currentKey string
		if keyPrefix == "" || keyPrefix == "-" {
			currentKey = key
		} else {
			currentKey = fmt.Sprintf("%s%s%s", keyPrefix, "_", key)
		}

		if kindOfValue == reflect.Struct {
			err := loadFromEnvToModel(currentKey, fieldValue.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}

		envValue, envExists := envReader.LookupEnv(currentKey)

		if field.isRequired() && !envExists {
			return fmt.Errorf("required field %s is not set", key)
		}

		if !envExists && kindOfValue != reflect.Struct {
			if defaultValue, ok := field.getDefaultValue(); ok {
				envValue = defaultValue
			}
		}

		if envValue == "" {
			continue
		}

		switch kindOfValue {
		case reflect.String:
			fieldValue.SetString(envValue)

		case reflect.Int:
			intValue, err := strconv.Atoi(envValue)
			if err != nil {
				return &ErrParseEnvValue{
					Key:   currentKey,
					Value: envValue,
				}
			}
			fieldValue.SetInt(int64(intValue))

		case reflect.Float64:
			floatValue, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return &ErrParseEnvValue{
					Key:   currentKey,
					Value: envValue,
				}
			}
			fieldValue.SetFloat(floatValue)

		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envValue)
			if err != nil {
				return &ErrParseEnvValue{
					Key:   currentKey,
					Value: envValue,
				}
			}
			fieldValue.SetBool(boolValue)

		case reflect.Slice:
			sliceValue := strings.Split(envValue, ",")
			fieldValue.Set(reflect.ValueOf(sliceValue))

		case reflect.Map:
			err := loadFromEnvToMap(envValue, fieldValue)
			if err != nil {
				return &ErrParseEnvValue{
					Key:   currentKey,
					Value: envValue,
				}
			}
		}
	}

	return nil

}

func loadFromEnv(model any) error {
	// check the model type
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("model must be a pointer")
	}

	if reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("model must be a pointer to a struct")
	}

	// find all env keys and set to model
	return loadFromEnvToModel("", model)
}

// Loads the environment variables into the provided model
func Load(model any) error {
	err := loadFromEnv(model)
	if err != nil {
		return err
	}

	return nil
}
