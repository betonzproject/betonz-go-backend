package product

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/doorman2137/betonz-go/internal/etg"
	"github.com/doorman2137/betonz-go/internal/utils/numericutils"
	"github.com/jackc/pgx/v5/pgtype"
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

func GetUserBalance(etgUsername string, product Product) (pgtype.Numeric, error) {
	payload := BalanceRequest{
		Op:   os.Getenv("ETG_OPERATOR_CODE"),
		Prod: product,
		Mem:  etgUsername,
		Pass: "00000000",
	}
	var balance pgtype.Numeric
	var balanceResponse BalanceResponse
	err := etg.Post("/balance", payload, &balanceResponse)
	if err != nil {
		return balance, err
	}

	if balanceResponse.Err != etg.Success {
		return balance, fmt.Errorf("%d: %s\nEndpoint: %s\nPayload: %+v", balanceResponse.Err, balanceResponse.Desc, "/balance", payload)
	}

	err = balance.Scan(strconv.FormatFloat(balanceResponse.Balance, 'f', 2, 64))

	return balance, err
}

type DepositRequest struct {
	Op     string         `json:"op"`
	Prod   Product        `json:"prod"`
	Mem    string         `json:"mem"`
	Pass   string         `json:"pass"`
	RefNo  string         `json:"ref_no"`
	Amount pgtype.Numeric `json:"amount"`
}

type DepositResponse struct {
	Err  int    `json:"err"`
	Desc string `json:"desc"`
}

func Deposit(refId string, etgUsername string, product Product, amount pgtype.Numeric) error {
	payload := DepositRequest{
		Op:     os.Getenv("ETG_OPERATOR_CODE"),
		Prod:   product,
		Mem:    etgUsername,
		Pass:   "00000000",
		RefNo:  refId,
		Amount: amount,
	}
	var depositResponse DepositResponse
	err := etg.Post("/deposit", payload, &depositResponse)
	if err != nil {
		return err
	}

	if depositResponse.Err != etg.Success {
		return fmt.Errorf("%d: %s\nEndpoint: %s\nPayload: %+v", depositResponse.Err, depositResponse.Desc, "/deposit", payload)
	}

	return nil
}

type WithdrawRequest struct {
	Op     string         `json:"op"`
	Prod   Product        `json:"prod"`
	Mem    string         `json:"mem"`
	Pass   string         `json:"pass"`
	Ref_No string         `json:"ref_no"`
	Amount pgtype.Numeric `json:"amount"`
}

type WithdrawResponse struct {
	Err  int    `json:"err"`
	Desc string `json:"desc"`
}

func Withdraw(refId string, etgUsername string, product Product, amount pgtype.Numeric) error {
	payload := WithdrawRequest{
		Op:     os.Getenv("ETG_OPERATOR_CODE"),
		Prod:   product,
		Mem:    etgUsername,
		Pass:   "00000000",
		Ref_No: refId,
		Amount: amount,
	}
	var withdrawResponse WithdrawResponse
	err := etg.Post("/withdraw", payload, &withdrawResponse)
	if err != nil {
		return err
	}

	if withdrawResponse.Err != etg.Success {
		return fmt.Errorf("%d: %s\nEndpoint: %s\nPayload: %+v", withdrawResponse.Err, withdrawResponse.Desc, "/withdraw", payload)
	}

	return nil
}

// Transfers from wallet `fromWallet` to `toWallet`.
func Transfer(q *db.Queries, ctx context.Context, refId string, user db.User, fromWallet Product, toWallet Product, amount pgtype.Numeric) error {
	if fromWallet != toWallet {
		if fromWallet == MainWallet {
			err := q.UpdateUserMainWallet(ctx, db.UpdateUserMainWalletParams{
				ID:         user.ID,
				MainWallet: numericutils.Sub(user.MainWallet, amount),
			})
			if err != nil {
				return err
			}

			return Deposit(refId, user.EtgUsername, toWallet, amount)
		}

		if toWallet == MainWallet {
			err := q.UpdateUserMainWallet(ctx, db.UpdateUserMainWalletParams{
				ID:         user.ID,
				MainWallet: numericutils.Add(user.MainWallet, amount),
			})
			if err != nil {
				return err
			}

			return Withdraw(refId, user.EtgUsername, fromWallet, amount)
		}

		err := Withdraw(refId, user.EtgUsername, fromWallet, amount)
		if err != nil {
			return err
		}

		err = Deposit(refId, user.EtgUsername, toWallet, amount)
		if err != nil {
			// Deposit to `toWallet` failed. Undo transfer by depositing back to `fromWallet`
			err2 := Deposit(refId, user.EtgUsername, fromWallet, amount)
			if err2 != nil {
				// Undo failed! Last resort deposit back to main wallet
				q.UpdateUserMainWallet(ctx, db.UpdateUserMainWalletParams{
					ID:         user.ID,
					MainWallet: numericutils.Add(user.MainWallet, amount),
				})
			}
			return err
		}
	}
	return nil
}
