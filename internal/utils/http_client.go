package utils

import (
	"encoding/json"
	"net/http"
)

func HttpGet[T any](url string, jsonResponse T) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(jsonResponse)
	return err
}
