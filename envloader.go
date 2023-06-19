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
	options   Options
	envReader EnvReader
}

type Options struct {
	EnvFieldDelimiter string
}

// DefaultOptions returns the default options
func DefaultOptions() Options {
	return Options{
		EnvFieldDelimiter: ".",
	}
}

// OptionFunc defines the signature of an option function
type OptionFunc func(*Options)

// WithEnvFieldDelimiter sets the delimiter used to separate the struct field names
func WithEnvFieldDelimiter(value string) OptionFunc {
	return func(opts *Options) {
		opts.EnvFieldDelimiter = value
	}
}

// New creates a new instance of envLoader with the provided options
// If no options are provided, the default options are used
// Default delimiter is "."
// You can change the delimiter with WithEnvFieldDelimiter option
// Example:
//
//	loader := envloader.New(envloader.WithEnvFieldDelimiter("*"))
func New(opts ...OptionFunc) EnvLoader {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &envLoader{
		options:   options,
		envReader: &DefaultEnvReader{},
	}
}

func (el *envLoader) toSnakeUpperCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
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
			currentKey = fmt.Sprintf("%s%s%s", keyPrefix, el.options.EnvFieldDelimiter, key)
		}

		envValue, exists := el.envReader.LookupEnv(currentKey)

		switch kindOfValue {
		case reflect.String:
			if exists {
				fieldValue.SetString(envValue)
			}

		case reflect.Int:
			if exists {
				intValue, err := strconv.Atoi(envValue)
				if err != nil {
					return err
				}
				fieldValue.SetInt(int64(intValue))
			}

		case reflect.Bool:
			if exists {
				boolValue, err := strconv.ParseBool(envValue)
				if err != nil {
					return err
				}
				fieldValue.SetBool(boolValue)
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
