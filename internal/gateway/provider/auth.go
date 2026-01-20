package provider

import (
	controllerAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/controller"
	repositoryAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/repository"
	routeAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/route"
	usecaseAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/usecase"

	"github.com/gin-gonic/gin"
)

var (
	authUserRepository              repositoryAuth.UserRepository
	authUserDetailRepository        repositoryAuth.UserDetailRepository
	authResetTokenRepository        repositoryAuth.ResetTokenRepository
	authVerificationTokenRepository repositoryAuth.VerificationTokenRepository
)

var (
	authAuthUsecase usecaseAuth.AuthUsecase
	profileUsecase  usecaseAuth.ProfileUsecase
	oauthUsecase    usecaseAuth.OauthUsecase
)

var (
	authAuthController *controllerAuth.AuthController
	profileController  *controllerAuth.ProfileController
	oauthController    *controllerAuth.OauthController
)

func ProvideAuthModule(router *gin.Engine) {
	injectAuthModuleRepository()
	injectAuthModuleUseCase()
	injectAuthModuleController()

	routeAuth.AuthControllerRoute(authAuthController, router, authMiddleware)
	routeAuth.ProfileControlRoute(profileController, router, authMiddleware)
	routeAuth.OauthControllerRoute(oauthController, router)
}

func injectAuthModuleRepository() {
	authUserRepository = repositoryAuth.NewUserRepository(dbWrapper)
	authUserDetailRepository = repositoryAuth.NewUserDetailRepository(dbWrapper, cfgConfig)
	authResetTokenRepository = repositoryAuth.NewResetTokenRepository(dbWrapper)
	authVerificationTokenRepository = repositoryAuth.NewVerificationTokenRepository(dbWrapper)
}

func injectAuthModuleUseCase() {
	authAuthUsecase = usecaseAuth.NewAuthUsecase(
		authUserRepository,
		authUserDetailRepository,
		redisUtil,
		jwtUtil,
		passwordEncryptor,
		base64Encryptor,
		emailTask,
		authResetTokenRepository,
		authVerificationTokenRepository,
		store,
	)
	profileUsecase = usecaseAuth.NewProfileUsecase(
		authUserRepository,
		authUserDetailRepository,
		redisUtil,
		passwordEncryptor,
		store,
	)
	oauthUsecase = usecaseAuth.NewOauthUsecase(authUserRepository, authUserDetailRepository, redisUtil, jwtUtil, store)
}

func injectAuthModuleController() {
	authAuthController = controllerAuth.NewAuthController(authAuthUsecase, store)
	profileController = controllerAuth.NewProfileController(profileUsecase, store)
	oauthController = controllerAuth.NewOauthController(oauthUsecase, cfgConfig)
}
