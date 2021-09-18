package config

import (
	"os"
	"reflect"
	"strconv"
)

func LoadConfigFromEnv() *Config {
	c := &Config{}
	v := reflect.ValueOf(c).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		envName := f.Tag.Get("env")
		defaultValue := f.Tag.Get("default")
		value := getEnv(envName, defaultValue)
		switch f.Type.Kind() {
		case reflect.Bool:
			cvalue, err := strconv.ParseBool(value)
			if err == nil {
				v.Field(i).SetBool(cvalue)
			}
		default:
			v.Field(i).SetString(value)
		}
	}
	return c
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
