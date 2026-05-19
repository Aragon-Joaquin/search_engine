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

	ENV_BOT_AGENT envVariable = "USERNAME_BOT_AGENT"
)

var all_env_vars = []envVariable{
	ENV_HOST, ENV_PORT, ENV_KEYHOST,
	ENV_VER, ENV_APPNAME, ENV_BOT_AGENT,
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	for _, e := range all_env_vars {
		setEnv(e)
	}
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
