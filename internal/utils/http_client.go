package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	client = http.Client{
		Timeout: 5 * time.Second,
	}
	agent_name string
)

func init() {
	agent_name = GetEnv(ENV_BOT_AGENT)
}

func HttpGet[T any](url string, jsonResponse T) error {
	if agent_name == "" {
		return fmt.Errorf("%s not found in .env (can get ip banned if is not set)", ENV_BOT_AGENT)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", agent_name)

	resp, err2 := client.Do(req)
	if err2 != nil {
		return err2
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(jsonResponse)
	return err
}
