package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/entity"
	apperrorPkg "github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	"github.com/faisalyudiansah/auth-service-template/pkg/database"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/transactor"

	"github.com/google/uuid"
)

type ResetTokenRepository interface {
	FindByResetToken(ctx context.Context, token, userID uuid.UUID) (*entity.ResetToken, error)
	Save(ctx context.Context, resetToken *entity.ResetToken) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type resetTokenRepositoryImpl struct {
	db database.Executor
}

func NewResetTokenRepository(db database.Executor) *resetTokenRepositoryImpl {
	return &resetTokenRepositoryImpl{
		db: db,
	}
}

func (r *resetTokenRepositoryImpl) FindByResetToken(ctx context.Context, token, userID uuid.UUID) (*entity.ResetToken, error) {
	query := `
		select id, created_at, updated_at from token_reset_users where reset_token = $1 and user_id = $2 and deleted_at is null
	`
	tx := transactor.ExtractTx(ctx)

	var (
		err        error
		resetToken = &entity.ResetToken{
			UserID:     userID,
			ResetToken: token,
		}
	)
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, token, userID).Scan(&resetToken.ID, &resetToken.CreatedAt, &resetToken.UpdatedAt)
	} else {
		err = r.db.QueryRowContext(ctx, query, token, userID).Scan(&resetToken.ID, &resetToken.CreatedAt, &resetToken.UpdatedAt)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrorPkg.NewNoRowsError(err, "reset token")
		}
		return nil, apperrorPkg.NewServerError(err)
	}

	return resetToken, nil
}

func (r *resetTokenRepositoryImpl) Save(ctx context.Context, resetToken *entity.ResetToken) error {
	query := `
		insert into token_reset_users(user_id, reset_token) values ($1, $2) returning id, created_at, updated_at
	`
	tx := transactor.ExtractTx(ctx)

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, resetToken.UserID, resetToken.ResetToken).Scan(&resetToken.ID, &resetToken.CreatedAt, &resetToken.UpdatedAt)
	} else {
		err = r.db.QueryRowContext(ctx, query, resetToken.UserID, resetToken.ResetToken).Scan(&resetToken.ID, &resetToken.CreatedAt, &resetToken.UpdatedAt)
	}

	if err != nil {
		return apperrorPkg.NewServerError(err)
	}

	return nil
}

func (r *resetTokenRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
		update token_reset_users set updated_at = now(), deleted_at = now() 
		where deleted_at is null and user_id = $1
	`
	tx := transactor.ExtractTx(ctx)

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, userID)
	} else {
		_, err = r.db.ExecContext(ctx, query, userID)
	}

	if err != nil {
		return apperrorPkg.NewServerError(err)
	}

	return nil
}
