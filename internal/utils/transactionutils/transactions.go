package transactionutils

import (
	"context"
	"log"

	"github.com/BetOnz-Company/betonz-go/internal/app"
	"github.com/BetOnz-Company/betonz-go/internal/db"

	"github.com/jackc/pgx/v5"
)

// Starts a database transaction. Panics if the transaction fails to start.
func Begin(app *app.App, ctx context.Context) (pgx.Tx, *db.Queries) {
	tx, err := app.Pool.Begin(ctx)
	if err != nil {
		log.Panicln("Can't start transaction: " + err.Error())
	}
	qtx := app.DB.WithTx(tx)
	return tx, qtx
}
