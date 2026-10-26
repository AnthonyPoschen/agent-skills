package config

import "os"

func Env() string {
	return os.Getenv("APP_ENV")
}
