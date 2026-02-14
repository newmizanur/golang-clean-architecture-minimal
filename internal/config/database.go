package config

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *bun.DB {
	username := viper.GetString("database.username")
	password := url.QueryEscape(viper.GetString("database.password"))
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	database := viper.GetString("database.name")
	idleConnection := viper.GetInt("database.pool.idle")
	maxConnection := viper.GetInt("database.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.pool.lifetime")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, database)

	sqldb, err := sql.Open("pg", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := sqldb.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	sqldb.SetMaxIdleConns(idleConnection)
	sqldb.SetMaxOpenConns(maxConnection)
	sqldb.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return bun.NewDB(sqldb, pgdialect.New())
}
