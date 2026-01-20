package transactor

import (
	"context"

	"github.com/faisalyudiansah/auth-service-template/pkg/database"
)

type Transactor interface {
	Atomic(ctx context.Context, fn func(context.Context) error) error
}

type transactorImpl struct {
	db *database.DB
}

func NewTransactor(db *database.DB) *transactorImpl {
	return &transactorImpl{
		db: db,
	}
}

func (t *transactorImpl) Atomic(ctx context.Context, fn func(context.Context) error) error {
	if existingTx := ExtractTx(ctx); existingTx != nil {
		return fn(ctx)
	}

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := injectTx(ctx, tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

type TxKey struct{}

func injectTx(ctx context.Context, tx *database.Tx) context.Context {
	return context.WithValue(ctx, TxKey{}, tx)
}

func ExtractTx(ctx context.Context) *database.Tx {
	if tx, ok := ctx.Value(TxKey{}).(*database.Tx); ok {
		return tx
	}
	return nil
}
