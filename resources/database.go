package resources

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"passVault/dtos"
	"time"
)

func initDatabaseConn() error {
	var (
		openDBOpts []stdlib.OptionOpenDB
		dsn        string
	)

	dsn = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.GetString(dtos.ConfigKeys.Database.Host),
		config.GetString(dtos.ConfigKeys.Database.Username),
		config.GetString(dtos.ConfigKeys.Database.Password),
		config.GetString(dtos.ConfigKeys.Database.Name),
		config.GetInt(dtos.ConfigKeys.Database.Port),
	)
	pgxConnConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return err
	}

	sqlDB := stdlib.OpenDB(*pgxConnConfig, openDBOpts...)

	sqlDB.SetMaxOpenConns(config.GetInt(dtos.ConfigKeys.Database.MaxOpenConnections))
	sqlDB.SetMaxIdleConns(config.GetInt(dtos.ConfigKeys.Database.MaxIdleConnections))
	sqlDB.SetConnMaxIdleTime(config.GetDuration(dtos.ConfigKeys.Database.MaxIdleConnectionTime) * time.Minute)

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		return err
	}

	db.Debug()

	databaseConn = db
	return nil
}
