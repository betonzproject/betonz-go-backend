package etg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const Success int = 1

func Post[T any](route string, payload any, dst *T) error {
	endpoint := os.Getenv("ETG_API_ENDPOINT") + route

	marshalled, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(marshalled))
	if err != nil {
		return fmt.Errorf("Can't create request: %s\nEndpoint: %s\nPayload: %+v", err, endpoint, payload)
	}
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + os.Getenv("AUTHORIZATION_TOKEN")},
		"Proxy-Url":     {os.Getenv("PROXY_URL")},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Can't read response body: %s\nEndpoint: %s\nPayload: %+v", err, endpoint, payload)
	}

	err = json.Unmarshal(body, dst)
	if err != nil {
		return fmt.Errorf("Can't unmarshal response body: %s\nEndpoint: %s\nPayload: %+v", err, endpoint, payload)
	}

	return nil
}
