package utils

import (
	"fmt"
	"os"
)

func MustGetEnvOrPanic(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(fmt.Errorf("environment variable '%s' is not set", k))
	}
	return v
}
