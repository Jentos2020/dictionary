package pg

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type PGStore struct {
	Schema string
	Db     *gorm.DB
	log    *slog.Logger
}

func New(
	dsn,
	dbSchema string,
	maxConn,
	maxIdle int,
	logs *slog.Logger,
) (*PGStore, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   dbSchema + ".",
			SingularTable: true,
		},
		Logger:         NewSlogAdapter(logs).LogMode(logger.Error),
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(maxConn)
	sqlDB.SetMaxIdleConns(maxIdle)

	ps := &PGStore{
		Schema: dbSchema,
		Db:     db,
		log:    logs,
	}

	return ps, nil
}

func (pgstore *PGStore) Ping(ctx context.Context) error {
	result, err := pgstore.Db.WithContext(ctx).DB()
	if err != nil {
		return errors.Wrap(err, "get db")
	}

	return result.PingContext(ctx)
}
