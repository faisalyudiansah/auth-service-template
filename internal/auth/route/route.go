package route

import (
	"github.com/faisalyudiansah/auth-service-template/internal/auth/controller"
	custom_type "github.com/faisalyudiansah/auth-service-template/internal/auth/entity/type"
	"github.com/faisalyudiansah/auth-service-template/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func AuthControllerRoute(c *controller.AuthController, r *gin.Engine, authMiddleware *middleware.AuthMiddleware) {
	g := r.Group("/auth")
	{
		g.POST("/login", c.Login)
		g.POST("/register", c.Register)
		g.POST("/register/from-admin", authMiddleware.Authorization(), authMiddleware.ProtectedRoles(custom_type.RoleAdmin), c.RegisterFromAdmin)
		g.POST("/refresh", authMiddleware.Authorization(), c.RefreshToken)
		g.POST("/logout", authMiddleware.Authorization(), c.Logout)
		g.POST("/forgot-password", c.ForgotPassword)
		g.POST("/reset-password", c.ResetPassword)
		g.POST("/send-verification", c.SendVerification)
		g.POST("/verify-account", c.VerifyAccount)
		g.PATCH("/inactive-account", authMiddleware.Authorization(), c.InactiveAccount)
	}
}

func ProfileControlRoute(c *controller.ProfileController, r *gin.Engine, authMiddleware *middleware.AuthMiddleware) {
	g := r.Group("/user", authMiddleware.Authorization())
	{
		g.GET("", authMiddleware.ProtectedRoles(custom_type.RoleAdmin), c.GetList)
		g.GET("/me", c.GetMe)
		g.GET("/:user_id", authMiddleware.ProtectedRoles(custom_type.RoleAdmin), c.GetUserByID)
		g.PUT("/:user_id", authMiddleware.OnlySelfOrAdmin("user_id"), c.UpdateUser)
		g.DELETE("/:user_id", authMiddleware.ProtectedRoles(custom_type.RoleAdmin), c.DeleteUser)
	}
}

func OauthControllerRoute(c *controller.OauthController, r *gin.Engine) {
	g := r.Group("/oauth")
	{
		g.GET("/:provider/login", c.Login)
		g.GET("/:provider/callback", c.Callback)
		g.GET("/:provider/logout", c.Logout)
	}
}
