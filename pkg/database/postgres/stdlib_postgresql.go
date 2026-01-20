package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/database"
	"github.com/faisalyudiansah/auth-service-template/pkg/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitStdLib(cfg *config.Config) *database.DB {
	dbCfg := cfg.Database

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
		dbCfg.Host,
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.DbName,
		dbCfg.Port,
		dbCfg.Sslmode,
	)

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log.Fatalf("error initializing database: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Log.Fatalf("error connecting to database: %v", err)
	}

	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(
		time.Duration(dbCfg.MaxConnLifetimeMinute) * time.Minute,
	)

	return &database.DB{
		DB:        sqlDB,
		Debug:     dbCfg.Debug,
		SlowLimit: time.Duration(dbCfg.SlowQueryMs) * time.Millisecond,
	}
}
