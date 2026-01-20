package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/faisalyudiansah/auth-service-template/internal/gateway/provider"
	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/logger"
	"github.com/faisalyudiansah/auth-service-template/pkg/middleware"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/validationutils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/shopspring/decimal"
	"golang.org/x/time/rate"
)

type HttpServer struct {
	cfg    *config.Config
	server *http.Server
}

func NewHttpServer(cfg *config.Config) *HttpServer {
	gin.SetMode(cfg.App.Environment)

	store := sessions.NewCookieStore([]byte(cfg.HttpServer.SessionSecret))
	store.MaxAge(86400 * cfg.HttpServer.SessionAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	if cfg.App.Environment == gin.ReleaseMode {
		store.Options.Secure = true
	}

	gothic.Store = store

	goth.UseProviders(
		google.New(cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.Google.CallbackURL),
	)

	router := gin.New()
	router.ContextWithFallback = true
	router.HandleMethodNotAllowed = true

	RegisterValidators()
	RegisterMiddleware(router, cfg)

	provider.ProvideHttpDependency(cfg, router)

	return &HttpServer{
		cfg: cfg,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port),
			Handler: router,
		},
	}
}

func (s *HttpServer) Start() {
	logger.Log.Info("Running HTTP server on port:", s.cfg.HttpServer.Port)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal("Error while HTTP server listening:", err)
	}
	logger.Log.Info("HTTP server is not receiving new requests...")
}

func (s *HttpServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.HttpServer.GracePeriod)*time.Second)
	defer cancel()

	logger.Log.Info("Attempting to shut down the HTTP server...")
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Error shutting down HTTP server:", err)
	}
	logger.Log.Info("HTTP server shut down gracefully")
}

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(validationutils.TagNameFormatter)
		v.RegisterCustomTypeFunc(validationutils.DecimalType, decimal.Decimal{})

		v.RegisterValidation("birth_date", validationutils.BirthDateValidator)
		v.RegisterValidation("url", validationutils.URLValidator)
		v.RegisterValidation("dgt", validationutils.DecimalGT)
		v.RegisterValidation("dlt", validationutils.DecimalLT)
		v.RegisterValidation("dgte", validationutils.DecimalGTE)
		v.RegisterValidation("dlte", validationutils.DecimalLTE)
		v.RegisterValidation("password", validationutils.PasswordValidator)
		v.RegisterValidation("clean_input", validationutils.CleanInputValidator)
		v.RegisterValidation("phone_number", validationutils.PhoneNumberValidator)
		v.RegisterValidation("time_format", validationutils.TimeFormatValidator)
		v.RegisterValidation("day_of_weeks", validationutils.DayOfWeekValidator)
		v.RegisterValidation("no_duplicates", validationutils.NoDuplicatesValidator)
	}
}

func RegisterMiddleware(router *gin.Engine, cfg *config.Config) {
	limiter := rate.NewLimiter(rate.Limit(cfg.HttpServer.MaxRequestPerSecond), cfg.HttpServer.MaxRequestPerSecond)

	middlewares := []gin.HandlerFunc{
		middleware.Logger(),
		middleware.Metrics(),
		middleware.ErrorHandler(),
		middleware.RateLimiter(limiter),
		middleware.RequestTimeout(cfg),
		cors.New(cors.Config{
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*", "Authorization", "Content-Type"},
			AllowOrigins: []string{
				cfg.URLClientConfig.URLCientVerificationEmail,
				cfg.URLClientConfig.URLClientForgotPassword,
			},
			AllowCredentials: true,
		}),
		gin.Recovery(),
	}

	router.Use(middlewares...)
}
