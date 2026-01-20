package redis

import (
	"fmt"

	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/logger"

	"github.com/RediSearch/redisearch-go/redisearch"
)

func InitRedisSearch(cfg *config.Config) *redisearch.Client {
	rds := redisearch.NewClient(fmt.Sprintf("%v:%v", cfg.Redis.Host, cfg.Redis.Port), cfg.Redis.SearchIndex)

	if cfg.App.Environment == "debug" || cfg.App.Environment == "development" {
		if err := rds.Drop(); err != nil {
			logger.Log.Infof("%v : no existing index...", err)
		}
	}
	logger.Log.Info("redisearch is ready...")

	return rds
}
