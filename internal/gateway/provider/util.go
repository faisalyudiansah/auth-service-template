package provider

import (
	"database/sql"
	"log"
	"time"

	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/transactor"
	"github.com/faisalyudiansah/auth-service-template/pkg/middleware"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/cloudinaryutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/encryptutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/jwtutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/redisutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/smtputils"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

var (
	cloudinaryUtil    cloudinaryutils.CloudinaryUtil
	jwtUtil           jwtutils.JwtUtilInterface
	smtpUtil          smtputils.SMTPUtils
	redisUtil         redisutils.RedisUtil
	passwordEncryptor encryptutils.PasswordEncryptor
	base64Encryptor   encryptutils.Base64Encryptor
	store             transactor.Transactor
	authMiddleware    *middleware.AuthMiddleware
	cronJob           *cron.Cron
)

func ProvideUtils(cfg *config.Config, db *sql.DB, rdb *redis.Client) {
	cloudinaryUtil = cloudinaryutils.NewCloudinaryUtil()
	jwtUtil = jwtutils.NewJwtUtil(cfg.Jwt)
	smtpUtil = smtputils.NewSMTPUtils(cfg.SMTP)
	passwordEncryptor = encryptutils.NewBcryptPasswordEncryptor(cfg.App.BCryptCost)
	base64Encryptor = encryptutils.NewBase64Encryptor()
	redisUtil = redisutils.NewRedisUtils(cfg.Redis, rdb)
	store = transactor.NewTransactor(dbWrapper)

	authMiddleware = middleware.NewAuthMiddleware(redisUtil, jwtUtil)

	wib, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Failed to load WIB timezone: %v", err)
	}
	cronJob = cron.New(cron.WithLocation(wib))
}
