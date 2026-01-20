package usecase

import (
	"context"
	"time"

	apperrorAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/apperror"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/entity"
	custom_typeAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/repository"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/utils"
	apperrorPkg "github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	constantPkg "github.com/faisalyudiansah/auth-service-template/pkg/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/transactor"
	custom_typePkg "github.com/faisalyudiansah/auth-service-template/pkg/entity/type"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/jwtutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/redisutils"

	"github.com/google/uuid"
	"github.com/markbates/goth"
)

type OauthUsecase interface {
	Login(ctx context.Context, request *goth.User) (*entity.User, *entity.Session, error)
}

type oauthUsecaseImpl struct {
	userRepo       repository.UserRepository
	userDetailRepo repository.UserDetailRepository
	redisUtil      redisutils.RedisUtil
	jwtUtil        jwtutils.JwtUtilInterface
	transactor     transactor.Transactor
}

func NewOauthUsecase(
	userRepo repository.UserRepository,
	userDetailRepo repository.UserDetailRepository,
	redisUtil redisutils.RedisUtil,
	jwtUtil jwtutils.JwtUtilInterface,
	transactor transactor.Transactor,
) *oauthUsecaseImpl {
	return &oauthUsecaseImpl{
		userRepo:       userRepo,
		userDetailRepo: userDetailRepo,
		redisUtil:      redisUtil,
		jwtUtil:        jwtUtil,
		transactor:     transactor,
	}
}

func (u *oauthUsecaseImpl) Login(ctx context.Context, request *goth.User) (*entity.User, *entity.Session, error) {
	currentTime := time.Now()
	recordUserDB, err := u.userRepo.Find(ctx, "email", request.Email)
	if err != nil {
		return nil, nil, apperrorPkg.NewServerError(err)
	}
	if recordUserDB != nil && !recordUserDB.IsOauth {
		return nil, nil, apperrorAuth.NewInvalidEmailAlreadyExists(err)
	}
	if recordUserDB == nil {
		recordUserDB = &entity.User{Email: request.Email}
		if err := u.userRepo.SaveOauth(ctx, recordUserDB); err != nil {
			return nil, nil, apperrorPkg.NewServerError(err)
		}

		birthDate, err := time.Parse(constantPkg.DEFAULT_DATE_ONLY, "2000-12-30")
		if err != nil {
			return nil, nil, apperrorPkg.NewServerError(err)
		}

		userDetail := new(entity.UserDetail)
		userDetail.UserID = recordUserDB.ID
		userDetail.FullName = request.Name
		userDetail.Sex = custom_typeAuth.SexOther
		userDetail.BirthDate = custom_typePkg.DateOnly(birthDate)
		userDetail.CreatedBy = recordUserDB.CreatedBy

		if _, err := u.userDetailRepo.Save(ctx, userDetail); err != nil {
			return nil, nil, apperrorPkg.NewServerError(err)
		}
	}

	jti := uuid.NewString()
	accessToken, err := u.jwtUtil.Sign(recordUserDB.ID, recordUserDB.Role, jti, currentTime)
	if err != nil {
		return nil, nil, apperrorPkg.NewServerError(err)
	}

	refreshToken, err := u.jwtUtil.SignRefresh(currentTime)
	if err != nil {
		return nil, nil, apperrorPkg.NewServerError(err)
	}

	sessionID := uuid.New()

	payloadSession := &entity.Session{
		UserID:       recordUserDB.ID,
		Role:         recordUserDB.Role,
		JTI:          jti,
		SessionID:    sessionID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		LoginAt:      uint64(currentTime.UnixMilli()),
	}

	err = u.redisUtil.Set(ctx, utils.SessionKey(sessionID), payloadSession, 24*time.Hour)
	if err != nil {
		return nil, nil, apperrorPkg.NewServerError(err)
	}

	return recordUserDB, payloadSession, nil
}
