package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/entity"
	custom_type "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	apperrorPkg "github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	"github.com/faisalyudiansah/auth-service-template/pkg/database"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/transactor"
	dtoPkg "github.com/faisalyudiansah/auth-service-template/pkg/dto"
	custom_typePkg "github.com/faisalyudiansah/auth-service-template/pkg/entity/type"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/queryBuilder"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetListUserWithDetail(ctx context.Context, req *dtoPkg.ListRequest) ([]*entity.User, error)
	GetTotalCount(ctx context.Context, req *dtoPkg.ListRequest) (uint64, error)
	Find(ctx context.Context, field string, value any) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) (*entity.User, error)
	SaveOauth(ctx context.Context, user *entity.User) error
	UpdatePassword(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
}

type userRepositoryImpl struct {
	db database.Executor
}

func NewUserRepository(db database.Executor) *userRepositoryImpl {
	return &userRepositoryImpl{
		db: db,
	}
}

var (
	userAllowedFilters = map[string]string{
		"email":        "u.email",
		"role":         "u.role",
		"is_verified":  "u.is_verified",
		"is_oauth":     "u.is_oauth",
		"is_active":    "u.is_active",
		"created_at":   "u.created_at",
		"full_name":    "ud.full_name",
		"sex":          "ud.sex",
		"phone_number": "ud.phone_number",
	}

	userAllowedSorts = map[string]string{
		"email":      "u.email",
		"role":       "u.role",
		"created_at": "u.created_at",
		"updated_at": "u.updated_at",
		"full_name":  "ud.full_name",
	}
)

func (r *userRepositoryImpl) GetListUserWithDetail(ctx context.Context, req *dtoPkg.ListRequest) ([]*entity.User, error) {
	db := r.db.QueryContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryContext
	}

	query := `
		SELECT
			u.id,
			u.role,
			u.email,
			u.is_verified,
			u.is_oauth,
			u.is_active,
			u.created_at,
			u.created_by,
			u.updated_at,
			u.updated_by,
			u.deleted_at,
			u.deleted_by,

			ud.id,
			ud.user_id,
			ud.full_name,
			ud.sex,
			ud.phone_number,
			ud.image_url,
			ud.birth_date,
			ud.created_at,
			ud.created_by,
			ud.updated_at,
			ud.updated_by,
			ud.deleted_at,
			ud.deleted_by
		FROM users u
		LEFT JOIN user_details ud 
			ON ud.user_id = u.id AND ud.deleted_at IS NULL
	`

	qb := queryBuilder.NewQueryBuilder(query).
		SetAllowedFilters(userAllowedFilters).
		SetAllowedSorts(userAllowedSorts).
		AddWhere("u.deleted_at IS NULL")

	if err := qb.ApplyFilters(req.GetFilters()); err != nil {
		msg := "invalid apply filters"
		return nil, apperrorPkg.NewClientError(err, &msg)
	}

	if req.HasSort() {
		if err := qb.ApplySort(req.GetSort()); err != nil {
			msg := "invalid apply sort"
			return nil, apperrorPkg.NewClientError(err, &msg)
		}
	} else {
		qb.AddDefaultSort("u.created_at", "DESC")
	}

	offset := (req.Page - 1) * req.Limit
	query, args := qb.BuildWithPagination(req.Limit, offset)

	rows, err := db(ctx, query, args...)
	if err != nil {
		return nil, apperrorPkg.NewServerError(err)
	}
	defer rows.Close()

	result := make([]*entity.User, 0)

	for rows.Next() {
		item := &entity.User{}

		var (
			udID        sql.NullString
			udUserID    sql.NullString
			fullName    sql.NullString
			sex         sql.NullInt16
			phoneNumber sql.NullString
			imageURL    sql.NullString
			birthDate   sql.NullTime
			udCreatedAt sql.NullTime
			udCreatedBy sql.NullString
			udUpdatedAt sql.NullTime
			udUpdatedBy sql.NullString
			udDeletedAt sql.NullTime
			udDeletedBy sql.NullString
		)

		err = rows.Scan(
			&item.ID,
			&item.Role,
			&item.Email,
			&item.IsVerified,
			&item.IsOauth,
			&item.IsActive,
			&item.CreatedAt,
			&item.CreatedBy,
			&item.UpdatedAt,
			&item.UpdatedBy,
			&item.DeletedAt,
			&item.DeletedBy,

			&udID,
			&udUserID,
			&fullName,
			&sex,
			&phoneNumber,
			&imageURL,
			&birthDate,
			&udCreatedAt,
			&udCreatedBy,
			&udUpdatedAt,
			&udUpdatedBy,
			&udDeletedAt,
			&udDeletedBy,
		)

		if err != nil {
			return nil, apperrorPkg.NewServerError(err)
		}

		if udID.Valid {
			item.UserDetail = &entity.UserDetail{}

			item.UserDetail.ID, _ = uuid.Parse(udID.String)
			item.UserDetail.UserID, _ = uuid.Parse(udUserID.String)
			item.UserDetail.FullName = fullName.String
			item.UserDetail.Sex = custom_type.Sex(sex.Int16)
			item.UserDetail.ImageURL = imageURL.String
			item.UserDetail.BirthDate = custom_typePkg.DateOnly(birthDate.Time)
			item.UserDetail.CreatedAt = udCreatedAt.Time
			item.UserDetail.CreatedBy, _ = uuid.Parse(udCreatedBy.String)

			if phoneNumber.Valid {
				item.UserDetail.PhoneNumber = &phoneNumber.String
			}

			if udUpdatedAt.Valid {
				item.UserDetail.UpdatedAt = &udUpdatedAt.Time
			}

			if udUpdatedBy.Valid {
				id, _ := uuid.Parse(udUpdatedBy.String)
				item.UserDetail.UpdatedBy = &id
			}

			if udDeletedAt.Valid {
				item.UserDetail.DeletedAt = &udDeletedAt.Time
			}

			if udDeletedBy.Valid {
				id, _ := uuid.Parse(udDeletedBy.String)
				item.UserDetail.DeletedBy = &id
			}

		} else {
			item.UserDetail = nil
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrorPkg.NewServerError(err)
	}

	return result, nil
}

func (r *userRepositoryImpl) GetTotalCount(ctx context.Context, req *dtoPkg.ListRequest) (uint64, error) {
	db := r.db.QueryRowContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryRowContext
	}

	query := `
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		LEFT JOIN user_details ud 
			ON ud.user_id = u.id AND ud.deleted_at IS NULL
	`

	qb := queryBuilder.NewQueryBuilder(query).
		SetAllowedFilters(userAllowedFilters).
		AddWhere("u.deleted_at IS NULL")

	if err := qb.ApplyFilters(req.GetFilters()); err != nil {
		msg := "invalid apply filters"
		return 0, apperrorPkg.NewClientError(err, &msg)
	}

	query, args := qb.Build()

	var count uint64
	if err := db(ctx, query, args...).Scan(&count); err != nil {
		return 0, apperrorPkg.NewServerError(err)
	}

	return count, nil
}

func (r *userRepositoryImpl) Find(ctx context.Context, field string, value any) (*entity.User, error) {
	db := r.db.QueryRowContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryRowContext
	}

	query := fmt.Sprintf(`
		select id, role, email, hash_password, is_verified, is_oauth, is_active, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
		from users 
		where %s = $1 and deleted_at is null
	`, field)

	user := &entity.User{}

	if err := db(ctx, query, value).Scan(
		&user.ID,
		&user.Role,
		&user.Email,
		&user.HashPassword,
		&user.IsVerified,
		&user.IsOauth,
		&user.IsActive,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.DeletedBy,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrorPkg.NewNoRowsError(err, value)
		}
		return nil, apperrorPkg.NewServerError(err)
	}
	return user, nil
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *entity.User) (*entity.User, error) {
	db := r.db.QueryRowContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryRowContext
	}

	query := `
		INSERT INTO users (id, role, email, hash_password, is_verified, created_at, created_by) VALUES 
		($1, $2, $3, $4, $5, NOW(), $6)
		RETURNING id, role, email, hash_password, is_verified, is_oauth, is_active, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by;
	`

	var resUser entity.User

	if err := db(ctx, query, user.ID, user.Role, user.Email, user.HashPassword, user.IsVerified, user.CreatedBy).Scan(
		&resUser.ID,
		&resUser.Role,
		&resUser.Email,
		&resUser.HashPassword,
		&resUser.IsVerified,
		&resUser.IsOauth,
		&resUser.IsActive,
		&resUser.CreatedAt,
		&resUser.CreatedBy,
		&resUser.UpdatedAt,
		&resUser.UpdatedBy,
		&resUser.DeletedAt,
		&resUser.DeletedBy,
	); err != nil {
		return nil, apperrorPkg.NewServerError(err)
	}

	return &resUser, nil
}

func (r *userRepositoryImpl) SaveOauth(ctx context.Context, user *entity.User) error {
	db := r.db.QueryRowContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryRowContext
	}

	id := uuid.New()

	query := `
		INSERT INTO users (id, role, email, is_verified, is_oauth, created_at, created_by) VALUES 
		($1, 1, $2, true, true, now(), $3)
		RETURNING id, role, email, hash_password, is_verified, is_oauth, is_active, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by;
	`

	if err := db(ctx, query, id, user.Email, id).Scan(
		&user.ID,
		&user.Role,
		&user.Email,
		&user.HashPassword,
		&user.IsVerified,
		&user.IsOauth,
		&user.IsActive,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.DeletedBy,
	); err != nil {
		return apperrorPkg.NewServerError(err)
	}

	return nil
}

func (r *userRepositoryImpl) UpdatePassword(ctx context.Context, user *entity.User) error {
	db := r.db.ExecContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.ExecContext
	}

	query := `
		update users set hash_password = $1, updated_at = now(), updated_by = $2 where id = $3
	`

	if _, err := db(ctx, query, user.HashPassword, user.UpdatedBy, user.ID); err != nil {
		return apperrorPkg.NewServerError(err)
	}

	return nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	db := r.db.QueryRowContext
	if tx := transactor.ExtractTx(ctx); tx != nil {
		db = tx.QueryRowContext
	}

	query := `
		UPDATE users
		SET
			role = $1,
			is_verified = $2,
			is_active = $3,
			updated_at = NOW(),
			updated_by = $4
	`

	args := []any{
		user.Role,
		user.IsVerified,
		user.IsActive,
		user.UpdatedBy,
	}

	argIdx := 5

	if user.DeletedBy != nil {
		query += `
			, deleted_at = NOW()
			, deleted_by = $` + strconv.Itoa(argIdx)
		args = append(args, *user.DeletedBy)
		argIdx++

		query += `
			, deleted_reason = $` + strconv.Itoa(argIdx)
		args = append(args, user.DeletedReason)
		argIdx++
	}

	query += `
		WHERE id = $` + strconv.Itoa(argIdx) + `
		  AND deleted_at IS NULL
		RETURNING
			updated_at
	`
	args = append(args, user.ID)

	err := db(ctx, query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrorPkg.NewNoRowsError(err, user.ID)
		}
		return apperrorPkg.NewServerError(err)
	}

	return nil
}
