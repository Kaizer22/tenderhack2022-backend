package utils

import (
	"os"
	"strconv"
	"time"
	"unicode/utf8"
)

//TODO to structure and move to a library

func StringInArray(val string, arr []string) (index int, contains bool) {
	for i, s := range arr {
		if s == val {
			return i, true
		}
	}
	return -1, false
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, dflt int64) int64 {
	if s, ok := os.LookupEnv(key); ok {
		if result, err := strconv.ParseInt(s, 10, 64); err != nil {
			return dflt
		} else {
			return result
		}
	}
	return dflt
}

func GetEnvBool(key string, dflt bool) bool {
	if s, ok := os.LookupEnv(key); ok {
		if result, err := strconv.ParseBool(s); err != nil {
			return dflt
		} else {
			return result
		}
	}
	return dflt
}

func GetEnvDuration(key string, dflt time.Duration) time.Duration {
	if s, ok := os.LookupEnv(key); ok {
		if result, err := time.ParseDuration(s); err != nil {
			return dflt
		} else {
			return result
		}
	}
	return dflt
}
