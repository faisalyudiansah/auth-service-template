package controller

import (
	"net/http"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/usecase"
	"github.com/faisalyudiansah/auth-service-template/internal/auth/utils"
	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/ginutils"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type OauthController struct {
	oauthUsecase usecase.OauthUsecase
	cfgConfig    *config.Config
}

func NewOauthController(
	oauthUsecase usecase.OauthUsecase,
	cfgConfig *config.Config,
) *OauthController {
	return &OauthController{
		oauthUsecase: oauthUsecase,
		cfgConfig:    cfgConfig,
	}
}

func (c *OauthController) Login(ctx *gin.Context) {
	provider := ctx.Param("provider")

	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (c *OauthController) Callback(ctx *gin.Context) {
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.Error(err)
		return
	}

	_, session, err := c.oauthUsecase.Login(ctx, &user)
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.SetSessionCookie(ctx, session.SessionID)

	ctx.Redirect(http.StatusFound, c.cfgConfig.URLClientConfig.URLClientOauthCallback)
}

func (c *OauthController) Logout(ctx *gin.Context) {
	provider := ctx.Param("provider")

	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()

	if err := gothic.Logout(ctx.Writer, ctx.Request); err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOKPlain(ctx)
}
