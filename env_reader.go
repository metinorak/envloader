package goenv

import "os"

type EnvReader interface {
	LookupEnv(key string) (string, bool)
}

type DefaultEnvReader struct{}

func (r *DefaultEnvReader) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}
