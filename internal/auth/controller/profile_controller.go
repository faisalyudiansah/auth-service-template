package controller

import (
	"fmt"

	"github.com/faisalyudiansah/auth-service-template/configs/logstash"
	apperrorAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/apperror"
	converterAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/dto/converter"
	dto_request "github.com/faisalyudiansah/auth-service-template/internal/auth/dto/request"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/usecase"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/utils"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/transactor"
	dtoPkg "github.com/faisalyudiansah/auth-service-template/pkg/dto"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/ginutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/pageutils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileController struct {
	profileUsecase usecase.ProfileUsecase
	transactor     transactor.Transactor
}

func NewProfileController(
	profileUsecase usecase.ProfileUsecase,
	transactor transactor.Transactor,
) *ProfileController {
	return &ProfileController{
		profileUsecase: profileUsecase,
		transactor:     transactor,
	}
}

func (c *ProfileController) GetList(ctx *gin.Context) {
	req := &dtoPkg.ListRequest{
		UserID: utils.GetValueUserIDFromContext(ctx),
	}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.Error(err)
		return
	}

	listUser, total, err := c.profileUsecase.GetList(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	res, paging := pageutils.CreateMetaData(
		converterAuth.ListUserEntityToDTOResponse(listUser),
		req.Page,
		req.Limit,
		total,
	)

	paging.Links = pageutils.CreateLinks(
		ctx.Request,
		int(paging.Page),
		int(paging.Size),
		int(paging.TotalItem),
		int(paging.TotalPage),
	)

	ginutils.ResponseOKPagination(ctx, res, paging)
}

func (c *ProfileController) GetMe(ctx *gin.Context) {
	res, err := c.profileUsecase.GetMe(ctx, utils.GetValueUserIDFromContext(ctx))
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOK(ctx, converterAuth.UserEntityToDTOResponse(res))
}

func (c *ProfileController) GetUserByID(ctx *gin.Context) {
	userIDstr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		ctx.Error(apperrorAuth.NewInvalidUserIdError())
		return
	}

	res, err := c.profileUsecase.GetUserByID(ctx, userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOK(ctx, converterAuth.UserEntityToDTOResponse(res))
}

func (c *ProfileController) UpdateUser(ctx *gin.Context) {
	modulName := "ProfileController.UpdateUser"

	req := new(dto_request.UpdateUser)
	userIDstr := ctx.Param("user_id")

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		logstash.LogstashError(ctx, err, userIDstr, fmt.Sprintf("%v - PARSE UUID", modulName))
		ctx.Error(apperrorAuth.NewInvalidUserIdError())
		return
	}

	if err := ctx.ShouldBindJSON(req); err != nil {
		logstash.LogstashError(ctx, err, nil, fmt.Sprintf("%v - BIND", modulName))
		ctx.Error(err)
		return
	}

	req.UserID = userID
	req.UpdatedBy = utils.GetValueUserIDFromContext(ctx)
	req.RoleWhoIsEdit = utils.GetValueRoleUserFromContext(ctx)

	logstash.LogstashRequestInfo(ctx, req, fmt.Sprintf("%v - REQUEST : %v", modulName, req.UserID))
	res, err := c.profileUsecase.UpdateUser(ctx, req)
	if err != nil {
		logstash.LogstashError(ctx, err, req, fmt.Sprintf("%v - USECASE : %s", modulName, req.UserID))
		ctx.Error(err)
		return
	}

	logstash.LogstashResponseInfo(ctx, res, fmt.Sprintf("%v - RESPONSE : %v", modulName, res.ID))
	ginutils.ResponseOK(ctx, converterAuth.UserEntityToDTOResponse(res))
}

func (c *ProfileController) DeleteUser(ctx *gin.Context) {
	modulName := "ProfileController.DeleteUser"

	req := new(dto_request.DeleteUser)
	userIDstr := ctx.Param("user_id")

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		logstash.LogstashError(ctx, err, userIDstr, fmt.Sprintf("%v - PARSE UUID", modulName))
		ctx.Error(apperrorAuth.NewInvalidUserIdError())
		return
	}

	if err := ctx.ShouldBindJSON(req); err != nil {
		logstash.LogstashError(ctx, err, nil, fmt.Sprintf("%v - BIND", modulName))
		ctx.Error(err)
		return
	}

	req.UserID = userID
	req.UpdatedBy = utils.GetValueUserIDFromContext(ctx)

	logstash.LogstashRequestInfo(ctx, req, fmt.Sprintf("%v - REQUEST : %v", modulName, req.UserID))
	err = c.profileUsecase.DeleteUser(ctx, req)
	if err != nil {
		logstash.LogstashError(ctx, err, req, fmt.Sprintf("%v - USECASE : %s", modulName, req.UserID))
		ctx.Error(err)
		return
	}

	ginutils.ResponseOKPlain(ctx)
}
