package usecase

import (
	"context"
	"errors"
	"time"

	apperrorAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/apperror"
	constantAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/constant"
	dto_request "github.com/faisalyudiansah/auth-service-template/internal/auth/dto/request"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/entity"
	custom_typeAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/repository"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/utils"
	"github.com/faisalyudiansah/auth-service-template/internal/queue/payload"
	"github.com/faisalyudiansah/auth-service-template/internal/queue/tasks"
	apperrorPkg "github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	constantPkg "github.com/faisalyudiansah/auth-service-template/pkg/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/transactor"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/encryptutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/jwtutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/redisutils"

	"github.com/google/uuid"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *dto_request.Login) (*entity.User, *entity.Session, error)
	RefreshToken(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	Register(ctx context.Context, req *dto_request.Register) (*entity.User, error)
	SendVerification(ctx context.Context, req *dto_request.SendVerification) error
	VerifyAccount(ctx context.Context, req *dto_request.VerifyAccount) error
	ForgotPassword(ctx context.Context, req *dto_request.ForgotPassword) error
	ResetPassword(ctx context.Context, req *dto_request.ResetPassword) error
	InactiveAccount(ctx context.Context, req *dto_request.InactiveAccount) error
}

type authUsecaseImpl struct {
	userRepo              repository.UserRepository
	userDetailRepo        repository.UserDetailRepository
	redisUtil             redisutils.RedisUtil
	jwtUtil               jwtutils.JwtUtilInterface
	passwordEncryptor     encryptutils.PasswordEncryptor
	base64Encryptor       encryptutils.Base64Encryptor
	emailTask             tasks.EmailTask
	resetTokenRepo        repository.ResetTokenRepository
	verificationTokenRepo repository.VerificationTokenRepository
	transactor            transactor.Transactor
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	userDetailRepo repository.UserDetailRepository,
	redisUtil redisutils.RedisUtil,
	jwtUtil jwtutils.JwtUtilInterface,
	passwordEncryptor encryptutils.PasswordEncryptor,
	base64Encryptor encryptutils.Base64Encryptor,
	emailTask tasks.EmailTask,
	resetTokenRepo repository.ResetTokenRepository,
	verificationTokenRepo repository.VerificationTokenRepository,
	transactor transactor.Transactor,
) *authUsecaseImpl {
	return &authUsecaseImpl{
		userRepo:              userRepo,
		userDetailRepo:        userDetailRepo,
		redisUtil:             redisUtil,
		jwtUtil:               jwtUtil,
		passwordEncryptor:     passwordEncryptor,
		base64Encryptor:       base64Encryptor,
		emailTask:             emailTask,
		resetTokenRepo:        resetTokenRepo,
		verificationTokenRepo: verificationTokenRepo,
		transactor:            transactor,
	}
}

func (u *authUsecaseImpl) Login(ctx context.Context, req *dto_request.Login) (*entity.User, *entity.Session, error) {
	currentTime := time.Now()

	recordUserDB, err := u.userRepo.Find(ctx, "email", req.Email)
	if err != nil {
		if err != apperrorPkg.NewNoRowsError(err, req.Email) {
			return nil, nil, apperrorAuth.NewEmailNotExistsError()
		}
		return nil, nil, err
	}

	if recordUserDB == nil || (recordUserDB.IsOauth && recordUserDB.HashPassword == "") {
		return nil, nil, apperrorAuth.NewEmailNotExistsError()
	}

	isValid := u.passwordEncryptor.Check(req.Password, recordUserDB.HashPassword)
	if !isValid {
		return nil, nil, apperrorAuth.NewInvalidLoginCredentials(nil)
	}

	if !recordUserDB.IsVerified {
		return nil, nil, apperrorAuth.NewUnverifiedError()
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

func (u *authUsecaseImpl) RefreshToken(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error) {
	currentTime := time.Now()

	getSession := &entity.Session{}
	sessionKey := utils.SessionKey(sessionID)

	err := u.redisUtil.GetWithScanJSON(ctx, sessionKey, getSession)
	if err != nil || getSession.UserID == uuid.Nil {
		return nil, apperrorPkg.NewForbiddenAccessError()
	}

	if getSession.RefreshToken == "" {
		return nil, apperrorPkg.NewForbiddenAccessError()
	}

	_, err = u.jwtUtil.Parse(getSession.RefreshToken)
	if err != nil {
		_ = u.redisUtil.Delete(ctx, sessionKey)
		return nil, apperrorPkg.NewSessionExpiredError()
	}

	user, err := u.userRepo.Find(ctx, "id", getSession.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperrorPkg.NewForbiddenAccessError()
	}

	newJTI := uuid.NewString()

	newAccessToken, err := u.jwtUtil.Sign(user.ID, user.Role, newJTI, currentTime)
	if err != nil {
		return nil, apperrorPkg.NewServerError(err)
	}

	newRefreshToken, err := u.jwtUtil.SignRefresh(currentTime)
	if err != nil {
		return nil, apperrorPkg.NewServerError(err)
	}

	updatedSession := &entity.Session{
		UserID:       user.ID,
		Role:         user.Role,
		JTI:          newJTI,
		SessionID:    sessionID,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		LoginAt:      getSession.LoginAt,
	}

	err = u.redisUtil.Set(ctx, sessionKey, updatedSession, 24*time.Hour)
	if err != nil {
		return nil, apperrorPkg.NewServerError(err)
	}

	return updatedSession, nil
}

func (u *authUsecaseImpl) Logout(ctx context.Context, sessionID uuid.UUID) error {
	err := u.redisUtil.Delete(ctx, utils.SessionKey(sessionID))
	if err != nil {
		return apperrorPkg.NewServerError(err)
	}

	return nil
}

func (u *authUsecaseImpl) Register(ctx context.Context, req *dto_request.Register) (*entity.User, error) {
	var entityUser *entity.User
	err := u.transactor.Atomic(ctx, func(txCtx context.Context) error {
		checkEmail, err := u.userRepo.Find(txCtx, "email", req.Email)
		if err != nil && err != apperrorPkg.NewNoRowsError(err, req.Email).OriginalError() {
			return err
		}
		if checkEmail != nil {
			return apperrorAuth.NewInvalidEmailAlreadyExists(nil)
		}
		hashPassword, err := u.passwordEncryptor.Hash(req.Password)
		if err != nil {
			return apperrorPkg.NewServerError(err)
		}

		user := new(entity.User)
		user.ID = uuid.New()
		user.Email = req.Email
		user.Role = custom_typeAuth.RoleUser
		user.HashPassword = hashPassword
		user.CreatedBy = user.ID

		if req.CreatedBy != nil { // admin create user
			user.Role = req.Role
			user.IsVerified = true
			user.CreatedBy = *req.CreatedBy
		}

		userDetail := new(entity.UserDetail)
		userDetail.UserID = user.ID
		userDetail.FullName = req.FullName
		userDetail.Sex = req.Sex
		userDetail.BirthDate = req.BirthDate
		userDetail.CreatedBy = user.CreatedBy

		resUser, err := u.userRepo.Save(txCtx, user)
		if err != nil {
			return err
		}

		resUserDetail, err := u.userDetailRepo.Save(txCtx, userDetail)
		if err != nil {
			return err
		}

		entityUser = resUser
		entityUser.UserDetail = resUserDetail
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entityUser, nil
}

func (u *authUsecaseImpl) SendVerification(ctx context.Context, req *dto_request.SendVerification) error {
	cachedToken := new(entity.VerificationToken)
	if err := u.redisUtil.GetWithScanJSON(ctx, utils.VerificationTokenCacheKey(req.Email), cachedToken); err == nil {
		if time.Since(cachedToken.CreatedAt) < constantAuth.VerificationTokenCooldownDuration {
			return apperrorAuth.NewTokenAlreadyExistsError()
		}
	}

	verificationTokenDb := new(entity.VerificationToken)
	err := u.transactor.Atomic(ctx, func(txCtx context.Context) error {
		userDb, err := u.userRepo.Find(txCtx, "email", req.Email)
		if err != nil {
			return err
		}
		if userDb == nil || userDb.IsOauth {
			return apperrorAuth.NewEmailNotExistsError()
		}
		if userDb.IsVerified {
			return apperrorAuth.NewVerifiedError()
		}

		userDetailDB, err := u.userDetailRepo.Find(txCtx, "user_id", userDb.ID)
		if err != nil {
			return err
		}
		if userDetailDB == nil {
			return apperrorAuth.NewUserDetailNotExistsError()
		}

		if err := u.verificationTokenRepo.DeleteByUserID(txCtx, userDb.ID); err != nil {
			return err
		}

		verificationTokenDb.UserID = userDb.ID
		verificationTokenDb.VerificationToken = uuid.New()
		if err := u.verificationTokenRepo.Save(txCtx, verificationTokenDb); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	if err := u.redisUtil.SetJSON(
		ctx,
		utils.VerificationTokenCacheKey(req.Email), verificationTokenDb, constantAuth.VerificationTokenCooldownDuration,
	); err != nil {
		return apperrorPkg.NewServerError(err)
	}

	if u.emailTask == nil {
		return apperrorPkg.NewServerError(errors.New("email task service not initialized"))
	}

	return u.emailTask.QueueVerificationEmail(ctx, &payload.VerificationEmailPayload{
		Email: req.Email,
		Token: verificationTokenDb.VerificationToken.String(),
	})
}

func (u *authUsecaseImpl) VerifyAccount(ctx context.Context, req *dto_request.VerifyAccount) error {
	email, err := u.base64Encryptor.DecodeURL(req.Email)
	if err != nil {
		return apperrorAuth.NewInvalidTokenCredentials()
	}

	verificationToken, err := u.base64Encryptor.DecodeURL(req.VerificationToken)
	if err != nil {
		return apperrorAuth.NewInvalidTokenCredentials()
	}

	parseVerificationToken, err := uuid.Parse(verificationToken)
	if err != nil {
		return apperrorAuth.NewInvalidTokenCredentials()
	}

	err = u.transactor.Atomic(ctx, func(txCtx context.Context) error {
		userDb, err := u.userRepo.Find(txCtx, "email", email)
		if err != nil {
			return err
		}
		if userDb == nil || userDb.IsOauth {
			return apperrorAuth.NewAccoundIsNotValidError()
		}
		if userDb.IsVerified {
			return apperrorAuth.NewVerifiedError()
		}

		userDetailDB, err := u.userDetailRepo.Find(txCtx, "user_id", userDb.ID)
		if err != nil {
			return err
		}
		if userDetailDB == nil {
			return apperrorAuth.NewUserDetailNotExistsError()
		}

		verificationTokenDb, err := u.verificationTokenRepo.FindByVerificationToken(txCtx, parseVerificationToken, userDb.ID)
		if err != nil && err != apperrorPkg.NewNoRowsError(err, "verification token").OriginalError() {
			return err
		}
		if verificationTokenDb == nil {
			return apperrorAuth.NewInvalidTokenCredentials()
		}
		if time.Now().After(verificationTokenDb.CreatedAt.Add(constantAuth.VerificationTokenExpireDuration - constantPkg.WIB)) {
			return apperrorAuth.NewExpiredTokenError()
		}

		userDb.IsVerified = true
		userDb.UpdatedBy = &userDb.ID
		if err := u.userRepo.Update(txCtx, userDb); err != nil {
			return err
		}
		if err := u.verificationTokenRepo.DeleteByUserID(txCtx, userDb.ID); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (u *authUsecaseImpl) ForgotPassword(ctx context.Context, req *dto_request.ForgotPassword) error {
	cachedToken := new(entity.ResetToken)
	if err := u.redisUtil.GetWithScanJSON(ctx, utils.ResetTokenCacheKey(req.Email), cachedToken); err == nil {
		if time.Since(cachedToken.CreatedAt) < constantAuth.ResetTokenCooldownDuration {
			return apperrorAuth.NewTokenAlreadyExistsError()
		}
	}

	resetTokenDb := new(entity.ResetToken)
	err := u.transactor.Atomic(ctx, func(txCtx context.Context) error {
		userDb, err := u.userRepo.Find(txCtx, "email", req.Email)
		if err != nil {
			return err
		}

		if userDb == nil || userDb.IsOauth {
			return apperrorAuth.NewEmailNotExistsError()
		}

		if err := u.resetTokenRepo.DeleteByUserID(txCtx, userDb.ID); err != nil {
			return err
		}

		resetTokenDb.UserID = userDb.ID
		resetTokenDb.ResetToken = uuid.New()
		if err := u.resetTokenRepo.Save(txCtx, resetTokenDb); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err := u.redisUtil.SetJSON(
		ctx,
		utils.ResetTokenCacheKey(req.Email), resetTokenDb, constantAuth.ResetTokenCooldownDuration,
	); err != nil {
		return apperrorPkg.NewServerError(err)
	}

	return u.emailTask.QueueForgotPasswordEmail(ctx, &payload.ForgotPasswordEmailPayload{
		Email: req.Email,
		Token: resetTokenDb.ResetToken.String(),
	})
}

func (u *authUsecaseImpl) ResetPassword(ctx context.Context, req *dto_request.ResetPassword) error {
	email, err := u.base64Encryptor.DecodeURL(req.Email)
	if err != nil {
		return apperrorAuth.NewInvalidTokenCredentials()
	}

	resetToken, err := u.base64Encryptor.DecodeURL(req.ResetToken)
	if err != nil {
		return apperrorAuth.NewInvalidTokenCredentials()
	}

	parseResetToken, err := uuid.Parse(resetToken)
	if err != nil {
		return apperrorAuth.NewInvalidTokenCredentials()
	}

	err = u.transactor.Atomic(ctx, func(txCtx context.Context) error {
		userDb, err := u.userRepo.Find(txCtx, "email", email)
		if err != nil {
			return err
		}
		if userDb == nil || userDb.IsOauth {
			return apperrorAuth.NewEmailNotExistsError()
		}

		resetTokenDb, err := u.resetTokenRepo.FindByResetToken(txCtx, parseResetToken, userDb.ID)
		if err != nil && err != apperrorPkg.NewNoRowsError(err, "reset token").OriginalError() {
			return err
		}
		if resetTokenDb == nil {
			return apperrorAuth.NewInvalidTokenCredentials()
		}
		if time.Now().After(resetTokenDb.CreatedAt.Add(constantAuth.ResetTokenExpireDuration - constantPkg.WIB)) {
			return apperrorAuth.NewExpiredTokenError()
		}

		hashPassword, err := u.passwordEncryptor.Hash(req.Password)
		if err != nil {
			return apperrorPkg.NewServerError(err)
		}

		userDb.HashPassword = hashPassword
		if err := u.userRepo.UpdatePassword(txCtx, userDb); err != nil {
			return err
		}
		if err := u.resetTokenRepo.DeleteByUserID(txCtx, userDb.ID); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (u *authUsecaseImpl) InactiveAccount(ctx context.Context, req *dto_request.InactiveAccount) error {
	err := u.transactor.Atomic(ctx, func(txCtx context.Context) error {
		userDb, err := u.userRepo.Find(txCtx, "id", req.UserID)
		if err != nil {
			return err
		}

		userDb.IsActive = false
		userDb.UpdatedBy = &req.UpdatedBy

		err = u.userRepo.Update(txCtx, userDb)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
