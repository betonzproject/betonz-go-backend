package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type BalanceRequest struct {
	Op   string  `json:"op"`
	Prod Product `json:"prod"`
	Mem  string  `json:"mem"`
	Pass string  `json:"pass"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
	Err     int     `json:"err"`
	Desc    string  `json:"desc"`
}

func GetUserBalance(etgUsername string, product Product) (float64, error) {
	endpoint := os.Getenv("ETG_API_ENDPOINT") + "/balance"

	payload := BalanceRequest{
		Op:   os.Getenv("ETG_OPERATOR_CODE"),
		Prod: product,
		Mem:  etgUsername,
		Pass: "00000000",
	}
	marshalled, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(marshalled))
	if err != nil {
		log.Panicf("Can't create request: %s\nEndpoint: %s\nPayload: %+v\n", err, endpoint, payload)
	}
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + os.Getenv("AUTHORIZATION_TOKEN")},
		"Proxy-Url":     {os.Getenv("PROXY_URL")},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Can't read response body: %s\nEndpoint: %s\nPayload: %+v\n", err, endpoint, payload)
	}

	var balanceResponse BalanceResponse
	err = json.Unmarshal(body, &balanceResponse)
	if err != nil {
		log.Panicf("Can't unmarshal response body: %s\nEndpoint: %s\nPayload: %+v\n", err, endpoint, payload)
	}

	if balanceResponse.Err != Success {
		return 0, fmt.Errorf("%d: %s\nEndpoint: %s\nPayload: %+v", balanceResponse.Err, balanceResponse.Desc, endpoint, payload)
	}

	return balanceResponse.Balance, nil
}
