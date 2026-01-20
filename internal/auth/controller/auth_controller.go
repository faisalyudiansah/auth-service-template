package controller

import (
	"context"
	"fmt"

	"github.com/faisalyudiansah/auth-service-template/configs/logstash"
	apperrorAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/apperror"
	converterAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/dto/converter"
	dto_request "github.com/faisalyudiansah/auth-service-template/internal/auth/dto/request"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/entity"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/usecase"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/utils"
	apperrorPkg "github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/transactor"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/ginutils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	authUsecase usecase.AuthUsecase
	transactor  transactor.Transactor
}

func NewAuthController(
	authUsecase usecase.AuthUsecase,
	transactor transactor.Transactor,
) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
		transactor:  transactor,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	req := new(dto_request.Login)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	res, session, err := c.authUsecase.Login(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.SetSessionCookie(ctx, session.SessionID)

	ginutils.ResponseOK(ctx, converterAuth.UserEntityToDTOLogin(res))
}

func (c *AuthController) Logout(ctx *gin.Context) {
	if err := c.authUsecase.Logout(ctx, utils.GetValueSessionIDFromContext(ctx)); err != nil {
		ctx.Error(err)
		return
	}

	utils.ClearSessionCookie(ctx)

	ginutils.ResponseOKPlain(ctx)
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	res, err := c.authUsecase.RefreshToken(ctx, utils.GetValueSessionIDFromContext(ctx))
	if err != nil {
		if appErr, ok := err.(*apperrorPkg.AppError); ok && appErr.Error() == apperrorPkg.NewSessionExpiredError().Error() {
			utils.ClearSessionCookie(ctx)
		}
		ctx.Error(err)
		return
	}

	utils.SetSessionCookie(ctx, res.SessionID)

	ginutils.ResponseOKPlain(ctx)
}

func (c *AuthController) Register(ctx *gin.Context) {
	req := new(dto_request.Register)

	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	var res *entity.User
	var err error
	errRes := c.transactor.Atomic(ctx, func(txCtx context.Context) error {
		res, err = c.authUsecase.Register(txCtx, req)
		if err != nil {
			return err
		}

		if !res.IsVerified {
			if err := c.authUsecase.SendVerification(txCtx, &dto_request.SendVerification{Email: res.Email}); err != nil {
				return err
			}
		}
		return nil
	})

	if errRes != nil {
		ctx.Error(errRes)
		return
	}

	ginutils.ResponseCreated(ctx, converterAuth.UserEntityToDTORegister(res))
}

func (c *AuthController) RegisterFromAdmin(ctx *gin.Context) {
	req := new(dto_request.Register)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	userIDAdmin := utils.GetValueUserIDFromContext(ctx)
	req.CreatedBy = &userIDAdmin

	res, err := c.authUsecase.Register(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseCreated(ctx, converterAuth.UserEntityToDTORegister(res))
}

func (c *AuthController) SendVerification(ctx *gin.Context) {
	req := new(dto_request.SendVerification)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	if err := c.authUsecase.SendVerification(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOKPlain(ctx)
}

func (c *AuthController) VerifyAccount(ctx *gin.Context) {
	req := new(dto_request.VerifyAccount)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	if err := c.authUsecase.VerifyAccount(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOKPlain(ctx)
}

func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	req := new(dto_request.ForgotPassword)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	if err := c.authUsecase.ForgotPassword(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseCreatedPlain(ctx)
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	req := new(dto_request.ResetPassword)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	if err := c.authUsecase.ResetPassword(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOKPlain(ctx)
}

func (c *AuthController) InactiveAccount(ctx *gin.Context) {
	modulName := "AuthController.InactiveAccount"

	req := new(dto_request.InactiveAccount)
	userIDstr := ctx.Param("user_id")

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		logstash.LogstashError(ctx, err, userIDstr, fmt.Sprintf("%v - PARSE UUID", modulName))
		ctx.Error(apperrorAuth.NewInvalidUserIdError())
		return
	}

	req.UserID = userID
	req.UpdatedBy = utils.GetValueUserIDFromContext(ctx)

	logstash.LogstashRequestInfo(ctx, req, fmt.Sprintf("%v - REQUEST : %s", modulName, req.UserID))
	if err := c.authUsecase.InactiveAccount(ctx, req); err != nil {
		logstash.LogstashError(ctx, err, req, fmt.Sprintf("%v - USECASE : %s", modulName+".InactiveAccount", req.UserID))
		ctx.Error(err)
		return
	}

	if req.UserID == req.UpdatedBy {
		if err := c.authUsecase.Logout(ctx, req.UserID); err != nil {
			logstash.LogstashError(ctx, err, req, fmt.Sprintf("%v - USECASE : %s", modulName+".Logout", req.UserID))
			ctx.Error(err)
			return
		}
	}

	utils.ClearSessionCookie(ctx)

	ginutils.ResponseOKPlain(ctx)
}
