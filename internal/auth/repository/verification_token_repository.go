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

type VerificationTokenRepository interface {
	FindByVerificationToken(ctx context.Context, token, userID uuid.UUID) (*entity.VerificationToken, error)
	Save(ctx context.Context, verificationToken *entity.VerificationToken) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type verificationTokenRepositoryImpl struct {
	db database.Executor
}

func NewVerificationTokenRepository(db database.Executor) *verificationTokenRepositoryImpl {
	return &verificationTokenRepositoryImpl{
		db: db,
	}
}

func (r *verificationTokenRepositoryImpl) FindByVerificationToken(ctx context.Context, token, userID uuid.UUID) (*entity.VerificationToken, error) {
	db := r.db.QueryRowContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryRowContext
	}

	query := `
		select id, created_at, updated_at from token_verification_users where verification_token = $1 and user_id = $2 and deleted_at is null
	`

	verificationToken := &entity.VerificationToken{
		UserID:            userID,
		VerificationToken: token,
	}

	if err := db(ctx, query, token, userID).Scan(&verificationToken.ID, &verificationToken.CreatedAt, &verificationToken.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrorPkg.NewNoRowsError(err, "verification token")
		}
		return nil, apperrorPkg.NewServerError(err)
	}

	return verificationToken, nil
}

func (r *verificationTokenRepositoryImpl) Save(ctx context.Context, verificationToken *entity.VerificationToken) error {
	db := r.db.QueryRowContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryRowContext
	}

	query := `
		insert into token_verification_users(user_id, verification_token) values($1, $2) returning id, created_at, updated_at
	`

	if err := db(ctx, query, verificationToken.UserID, verificationToken.VerificationToken).Scan(
		&verificationToken.ID,
		&verificationToken.CreatedAt,
		&verificationToken.UpdatedAt,
	); err != nil {
		return apperrorPkg.NewServerError(err)
	}

	return nil
}

func (r *verificationTokenRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	db := r.db.ExecContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.ExecContext
	}

	query := `
		update token_verification_users set updated_at = now(), deleted_at = now() 
		where deleted_at is null and user_id = $1
	`

	if _, err := db(ctx, query, userID); err != nil {
		return apperrorPkg.NewServerError(err)
	}

	return nil
}
