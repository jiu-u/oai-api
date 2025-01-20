package common

import (
	"fmt"
	"os"
	"strconv"
)

type EnvValueType interface {
	~string | ~int | ~bool
}

func GetEnv[T EnvValueType](key string, defaultValue T, parseFunc func(string) (T, error)) T {
	if key == "" || os.Getenv(key) == "" {
		return defaultValue
	}
	valueStr := os.Getenv(key)
	parsedValue, err := parseFunc(valueStr)
	if err != nil {
		fmt.Printf("failed to parse %s: %s, using default value: %v\n", key, err.Error(), defaultValue)
		return defaultValue
	}
	return parsedValue
}

// Helper 函数
func parseInt(value string) (int, error) {
	return strconv.Atoi(value)
}

func parseBool(value string) (bool, error) {
	return strconv.ParseBool(value)
}
