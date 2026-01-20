package usecase

import (
	"context"

	apperrorAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/apperror"
	dto_request "github.com/faisalyudiansah/auth-service-template/internal/auth/dto/request"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/entity"
	custom_type "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/repository"
	apperrorPkg "github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/transactor"
	dtoPkg "github.com/faisalyudiansah/auth-service-template/pkg/dto"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/encryptutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/redisutils"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type ProfileUsecase interface {
	GetList(ctx context.Context, req *dtoPkg.ListRequest) ([]*entity.User, uint64, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	UpdateUser(ctx context.Context, req *dto_request.UpdateUser) (*entity.User, error)
	DeleteUser(ctx context.Context, req *dto_request.DeleteUser) error
}

type profileUsecaseImpl struct {
	userRepo          repository.UserRepository
	userDetailRepo    repository.UserDetailRepository
	redisUtil         redisutils.RedisUtil
	passwordEncryptor encryptutils.PasswordEncryptor
	transactor        transactor.Transactor
}

func NewProfileUsecase(
	userRepo repository.UserRepository,
	userDetailRepo repository.UserDetailRepository,
	redisUtil redisutils.RedisUtil,
	passwordEncryptor encryptutils.PasswordEncryptor,
	transactor transactor.Transactor,
) *profileUsecaseImpl {
	return &profileUsecaseImpl{
		userRepo:          userRepo,
		userDetailRepo:    userDetailRepo,
		redisUtil:         redisUtil,
		passwordEncryptor: passwordEncryptor,
		transactor:        transactor,
	}
}

func (u *profileUsecaseImpl) GetList(ctx context.Context, req *dtoPkg.ListRequest) ([]*entity.User, uint64, error) {
	var (
		recordUserDetailDB []*entity.User
		total              uint64
	)

	if err := req.DecodeFilters(); err != nil {
		msg := "invalid decode filters"
		return nil, 0, apperrorPkg.NewClientError(err, &msg)
	}

	if err := req.DecodeSort(); err != nil {
		msg := "invalid decode sort"
		return nil, 0, apperrorPkg.NewClientError(err, &msg)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		recordUserDetailDB, err = u.userRepo.GetListUserWithDetail(ctx, req)
		return err
	})

	g.Go(func() error {
		var err error
		total, err = u.userRepo.GetTotalCount(ctx, req)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	return recordUserDetailDB, total, nil
}

func (u *profileUsecaseImpl) GetMe(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	recordUserDB, err := u.userRepo.Find(ctx, "id", userID)
	if err != nil {
		return nil, err
	}

	recordUserDetailDB, err := u.userDetailRepo.Find(ctx, "user_id", userID)
	if err != nil {
		return nil, err
	}

	recordUserDB.UserDetail = recordUserDetailDB

	return recordUserDB, nil
}

func (u *profileUsecaseImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	res := &entity.User{}

	err := u.transactor.Atomic(ctx, func(cForTx context.Context) error {
		recordUserDB, err := u.userRepo.Find(cForTx, "id", userID)
		if err != nil {
			return err
		}

		recordUserDetailDB, err := u.userDetailRepo.Find(cForTx, "user_id", userID)
		if err != nil {
			return err
		}

		res = recordUserDB
		res.UserDetail = recordUserDetailDB

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *profileUsecaseImpl) UpdateUser(ctx context.Context, req *dto_request.UpdateUser) (*entity.User, error) {
	res := &entity.User{}

	err := u.transactor.Atomic(ctx, func(cForTx context.Context) error {
		recordUserDB, err := u.userRepo.Find(cForTx, "id", req.UserID)
		if err != nil {
			return err
		}

		recordUserDetailDB, err := u.userDetailRepo.Find(cForTx, "user_id", req.UserID)
		if err != nil {
			return err
		}

		isPhoneNumberAlreadyExists, err := u.userDetailRepo.Find(cForTx, "phone_number", req.PhoneNumber)
		if err != nil && err != apperrorPkg.NewNoRowsError(err, req.PhoneNumber).OriginalError() {
			return err
		}

		if isPhoneNumberAlreadyExists != nil && isPhoneNumberAlreadyExists.UserID != req.UserID {
			return apperrorAuth.NewInvalidPhoneNumberAlreadyExists()
		}

		if recordUserDB.Role.IsRoleAdmin() && req.RoleWhoIsEdit.IsRoleAdmin() && recordUserDB.ID != req.UpdatedBy {
			return apperrorPkg.NewDontHavePermissionErrorMessageError()
		}

		if req.RoleWhoIsEdit.IsRoleAdmin() {
			recordUserDB.Role = custom_type.Role(req.Role)
			recordUserDB.IsVerified = req.IsVerified
			recordUserDB.IsActive = req.IsActive
			recordUserDB.UpdatedBy = &req.UpdatedBy

			if err := u.userRepo.Update(cForTx, recordUserDB); err != nil {
				return err
			}
		}

		recordUserDetailDB.FullName = req.FullName
		recordUserDetailDB.Sex = custom_type.Sex(req.Sex)
		recordUserDetailDB.PhoneNumber = &req.PhoneNumber
		recordUserDetailDB.ImageURL = req.ImageURL
		recordUserDetailDB.UpdatedBy = &req.UpdatedBy

		if err := u.userDetailRepo.Update(cForTx, recordUserDetailDB); err != nil {
			return err
		}

		res = recordUserDB
		recordUserDB.UserDetail = recordUserDetailDB

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *profileUsecaseImpl) DeleteUser(ctx context.Context, req *dto_request.DeleteUser) error {
	err := u.transactor.Atomic(ctx, func(cForTx context.Context) error {
		recordUserDB, err := u.userRepo.Find(cForTx, "id", req.UserID)
		if err != nil {
			return err
		}

		recordUserDetailDB, err := u.userDetailRepo.Find(cForTx, "user_id", req.UserID)
		if err != nil {
			return err
		}

		if recordUserDB.Role.IsRoleAdmin() || recordUserDB.ID == req.UpdatedBy {
			return apperrorPkg.NewDontHavePermissionErrorMessageError()
		}

		recordUserDB.UpdatedBy = &req.UpdatedBy
		recordUserDB.DeletedBy = &req.UpdatedBy
		recordUserDB.DeletedReason = &req.DeletedReason

		recordUserDetailDB.UpdatedBy = &req.UpdatedBy
		recordUserDetailDB.DeletedBy = &req.UpdatedBy
		recordUserDetailDB.DeletedReason = &req.DeletedReason

		if err := u.userRepo.Update(cForTx, recordUserDB); err != nil {
			return err
		}

		if err := u.userDetailRepo.Update(cForTx, recordUserDetailDB); err != nil {
			return err
		}

		return nil
	})

	return err
}
