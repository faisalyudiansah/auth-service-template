package provider

import (
	"database/sql"
	"fmt"

	"github.com/faisalyudiansah/auth-service-template/configs/logstash"
	gatewayController "github.com/faisalyudiansah/auth-service-template/internal/gateway/controller"

	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/database"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/postgres"
	"github.com/faisalyudiansah/auth-service-template/pkg/database/redis"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	redisV9 "github.com/redis/go-redis/v9"
)

var (
	db          *sql.DB
	dbWrapper   *database.DB
	rdb         *redisV9.Client
	rds         *redisearch.Client
	asynqClient *asynq.Client
	cfgConfig   *config.Config
)

func InitProvider(cfg *config.Config) {
	cfgConfig = cfg

	dbWrapper = postgres.InitStdLib(cfgConfig)
	db = dbWrapper.DB

	logstash.InitLogstash(cfgConfig)

	rdb = redis.InitRedis(cfgConfig.Redis)
	rds = redis.InitRedisSearch(cfgConfig)
	asynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%v:%v", cfgConfig.Redis.Host, cfgConfig.Redis.Port),
	})

	ProvideUtils(cfgConfig, db, rdb)
}

func ProvideHttpDependency(cfg *config.Config, router *gin.Engine) {
	if emailTask == nil {
		injectQueueModuleTask(asynqClient)
	}
	ProvideGatewayModule(router)
	ProvideAuthModule(router)
	cronJob.Start()
}

func ProvideQueueDependency(client *asynq.Client, mux *asynq.ServeMux) {
	ProvideQueueModule(client, mux, cfgConfig)
}

func ProvideGatewayModule(router *gin.Engine) {
	appController := gatewayController.NewAppController(db, rdb)
	appController.Route(router)
}
