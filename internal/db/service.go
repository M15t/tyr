package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/imdatngo/gowhere"

	"tyr/config"

	_ "gorm.io/driver/postgres" // DB adapter
	"gorm.io/gorm"

	// EnablePostgreSQL: remove the mysql package above, uncomment the following

	dbutil "github.com/M15t/gram/pkg/util/db"
	dblogger "github.com/M15t/gram/pkg/util/db/logger"
	"github.com/M15t/gram/pkg/util/prettylog"
)

// New creates new database connection to the database server
func New(cfg config.DB) (*gorm.DB, *sql.DB, error) {
	// Add your DB related stuffs here, such as:
	// - gorm.DefaultTableNameHandler
	// - gowhere.DefaultConfig

	gowhere.DefaultConfig.Dialect = gowhere.DialectPostgreSQL

	// logger config
	var levelLog slog.Leveler
	switch cfg.Logging {
	case -4:
		levelLog = slog.LevelDebug
	case 0:
		levelLog = slog.LevelInfo
	case 4:
		levelLog = slog.LevelWarn
	case 8:
		levelLog = slog.LevelError
	}

	var slogger *slog.Logger
	slogger = slog.New(prettylog.NewHandler(&slog.HandlerOptions{
		Level: levelLog,
	}, prettylog.JSONFormat))

	if config.IsLambda() {
		slogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: levelLog,
		}))
	}

	gcfg := dblogger.NewConfig(slogger.Handler()).
		WithTraceAll(true).
		WithRequestID(true)

	myLogger := dblogger.NewWithConfig(gcfg)

	// parse extra params, merge with default params
	// change to PostgreSQLDialector{} for PostgreSQL
	dbConn, err := dbutil.NewDBConnection(&dbutil.PostgreSQLDialector{}, dbutil.Config{
		Username: cfg.Username,
		Password: cfg.Password,
		Host:     cfg.Host,
		Port:     cfg.Port,
		Database: cfg.Database,
		Params:   cfg.Params,
	})
	if err != nil {
		return nil, nil, err
	}

	// connect to db server
	db, err := gorm.Open(dbConn, &gorm.Config{
		Logger:                                   myLogger,
		AllowGlobalUpdate:                        false,
		CreateBatchSize:                          1000,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("cannot establish connection: %w", err)
	}

	// connection pool settings
	sqldb, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get generic db instance: %w", err)
	}
	//! NOTE: These are not one-size-fits-all settings. Turn it based on your db settings!
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(10)
	sqldb.SetConnMaxLifetime(30 * time.Minute)
	sqldb.SetConnMaxIdleTime(10 * time.Minute)

	return db, sqldb, nil
}
