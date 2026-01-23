package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/entity"
	apperrorPkg "github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/database"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/transactor"
)

type UserDetailRepository interface {
	Find(ctx context.Context, field string, value any) (*entity.UserDetail, error)
	Save(ctx context.Context, userDetail *entity.UserDetail) (*entity.UserDetail, error)
	Update(ctx context.Context, userDetail *entity.UserDetail) error
}

type userDetailRepositoryImpl struct {
	db        database.Executor
	cfgConfig *config.Config
}

func NewUserDetailRepository(db database.Executor, cfgConfig *config.Config) *userDetailRepositoryImpl {
	return &userDetailRepositoryImpl{
		db:        db,
		cfgConfig: cfgConfig,
	}
}

func (r *userDetailRepositoryImpl) Find(ctx context.Context, field string, value any) (*entity.UserDetail, error) {
	db := r.db.QueryRowContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryRowContext
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, full_name, sex, phone_number, image_url, birth_date, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM user_details
		WHERE %s = $1 AND deleted_at IS NULL
	`, field)

	var userDetail entity.UserDetail

	if err := db(ctx, query, value).Scan(
		&userDetail.ID,
		&userDetail.UserID,
		&userDetail.FullName,
		&userDetail.Sex,
		&userDetail.PhoneNumber,
		&userDetail.ImageURL,
		&userDetail.BirthDate,
		&userDetail.CreatedAt,
		&userDetail.CreatedBy,
		&userDetail.UpdatedAt,
		&userDetail.UpdatedBy,
		&userDetail.DeletedAt,
		&userDetail.DeletedBy,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrorPkg.NewNoRowsError(err, value)
		}
		return nil, apperrorPkg.NewServerError(err)
	}

	return &userDetail, nil
}

func (r *userDetailRepositoryImpl) Save(ctx context.Context, userDetail *entity.UserDetail) (*entity.UserDetail, error) {
	db := r.db.QueryRowContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryRowContext

	}

	query := `
		INSERT INTO user_details (user_id, full_name, sex, phone_number, image_url, birth_date, created_at, created_by) VALUES 
		($1, $2, $3, NULL, $4, $5, NOW(), $6)
		RETURNING id, user_id, full_name, sex, phone_number, image_url, birth_date, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by;
	`
	defaultImgUrl := r.cfgConfig.App.DefaultImageUserProfile

	var resUserDetail entity.UserDetail

	if err := db(ctx, query, userDetail.UserID, userDetail.FullName, userDetail.Sex, defaultImgUrl, userDetail.BirthDate, userDetail.CreatedBy).Scan(
		&resUserDetail.ID,
		&resUserDetail.UserID,
		&resUserDetail.FullName,
		&resUserDetail.Sex,
		&resUserDetail.PhoneNumber,
		&resUserDetail.ImageURL,
		&resUserDetail.BirthDate,
		&resUserDetail.CreatedAt,
		&resUserDetail.CreatedBy,
		&resUserDetail.UpdatedAt,
		&resUserDetail.UpdatedBy,
		&resUserDetail.DeletedAt,
		&resUserDetail.DeletedBy,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrorPkg.NewNoRowsError(err, userDetail.UserID)
		}
		return nil, apperrorPkg.NewServerError(err)
	}

	return &resUserDetail, nil
}

func (r *userDetailRepositoryImpl) Update(ctx context.Context, userDetail *entity.UserDetail) error {
	db := r.db.QueryRowContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryRowContext
	}

	query := `
		UPDATE user_details
		SET
			full_name = $1,
			sex = $2,
			phone_number = $3,
			image_url = $4,
			birth_date = $5
	`

	args := []any{
		userDetail.FullName,
		userDetail.Sex,
		userDetail.PhoneNumber,
		userDetail.ImageURL,
		userDetail.BirthDate,
	}

	argIdx := 6

	if userDetail.UpdatedBy != nil && userDetail.DeletedBy == nil {
		query += `
			, updated_at = NOW()
			, updated_by = $` + strconv.Itoa(argIdx)
		args = append(args, *userDetail.UpdatedBy)
		argIdx++
	}

	if userDetail.DeletedBy != nil {
		query += `
			, deleted_at = NOW()
			, deleted_by = $` + strconv.Itoa(argIdx)
		args = append(args, *userDetail.DeletedBy)
		argIdx++

		query += `
			, deleted_reason = $` + strconv.Itoa(argIdx)
		args = append(args, userDetail.DeletedReason)
		argIdx++
	}

	query += `
		WHERE user_id = $` + strconv.Itoa(argIdx) + `
		  AND deleted_at IS NULL
		RETURNING
			updated_at
	`
	args = append(args, userDetail.UserID)

	err := db(ctx, query, args...).Scan(&userDetail.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrorPkg.NewNoRowsError(err, userDetail.ID)
		}
		return apperrorPkg.NewServerError(err)
	}

	return nil
}
