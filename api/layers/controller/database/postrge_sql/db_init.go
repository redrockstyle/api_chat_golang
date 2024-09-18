package postrge_sql

import (
	config "api_chat/api/layers/domain/cfg"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPsqlDatabase(cfg *config.Configuration) (db *gorm.DB, err error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Postgres.PostgresqlHost,
		cfg.Postgres.PostgresqlPort,
		cfg.Postgres.PostgresqlUser,
		cfg.Postgres.PostgresqlDB,
		cfg.Postgres.PostgresqlPass,
	)
	gormDB, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
