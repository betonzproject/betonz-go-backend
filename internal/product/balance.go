package product

import (
	"fmt"
	"os"

	"github.com/doorman2137/betonz-go/internal/etg"
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
	payload := BalanceRequest{
		Op:   os.Getenv("ETG_OPERATOR_CODE"),
		Prod: product,
		Mem:  etgUsername,
		Pass: "00000000",
	}
	var balanceResponse BalanceResponse
	err := etg.Post("/balance", payload, &balanceResponse)
	if err != nil {
		return 0, err
	}

	if balanceResponse.Err != etg.Success {
		return 0, fmt.Errorf("%d: %s\nEndpoint: %s\nPayload: %+v", balanceResponse.Err, balanceResponse.Desc, "/balance", payload)
	}

	return balanceResponse.Balance, nil
}
