package dao

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/jackc/pgx/v5"
)

// A TxFn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(tx pgx.Tx) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn` or when it panics.
func WithTransaction(ctx context.Context, fn TxFn) error {
	logger := ctxval.Logger(ctx)
	tx, beginErr := db.Pool.Begin(ctx)
	if beginErr != nil {
		logger.Warn().Err(beginErr).Msg("Cannot begin database transaction")
		return fmt.Errorf("transaction error: %w", beginErr)
	}

	defer func() {
		if p := recover(); p != nil {
			logger.Warn().Msgf("Rolling database transaction back due to panic call: %s", p)
			rollErr := tx.Rollback(ctx)
			if rollErr != nil {
				logger.Warn().Err(rollErr).Msg("Cannot rollback database transaction")
				return
			}
			panic(p)
		}
	}()

	callErr := fn(tx)

	if callErr != nil {
		logger.Warn().Msg("Rolling database transaction back due to error")
		rollErr := tx.Rollback(ctx)
		if rollErr != nil {
			logger.Warn().Err(rollErr).Msg("Cannot rollback database transaction")
			// return the call (root cause) error and not transaction error
			return fmt.Errorf("transaction error: %w", callErr)
		}
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		logger.Warn().Err(commitErr).Msg("Cannot rollback database transaction")
		return fmt.Errorf("transaction error: %w", commitErr)
	}

	return nil
}
