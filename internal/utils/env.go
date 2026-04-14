package utils

import (
	"os"

	"github.com/joho/godotenv"
)

// i know this is over-engineering. trying something new

type envValue struct {
	Value string
	Ok    bool
}

var envVariables = map[envVariable]envValue{}

type envVariable string

const (
	ENV_HOST    envVariable = "HOST"
	ENV_PORT    envVariable = "PORT"
	ENV_KEYHOST envVariable = "KEYHOST_PATH"
	ENV_VER     envVariable = "VERSION"
	ENV_APPNAME envVariable = "APP_NAME"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	setEnv(ENV_HOST)
	setEnv(ENV_PORT)
	setEnv(ENV_KEYHOST)
	setEnv(ENV_VER)
	setEnv(ENV_APPNAME)
}

func setEnv(id envVariable) {
	val, ok := os.LookupEnv(string(id))

	envVariables[id] = envValue{
		Value: val,
		Ok:    ok,
	}
}

func GetEnv(id envVariable) string {
	if val, ok := envVariables[id]; ok && val.Ok {
		return val.Value
	}
	return ""
}
