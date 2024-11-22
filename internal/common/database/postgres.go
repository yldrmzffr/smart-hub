package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"smart-hub/internal/common/logger"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

type PostgreConfig struct {
	DSN string
}

func NewPostgresDatabase(ctx context.Context, config *PostgreConfig) (Database, error) {
	logger.Info("PostgreSQL Starting...")

	poolConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		logger.Error("Error parsing PostgreSQL config: ", err)
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Error("PostgreSQL Connection Error: ", err)
		return nil, err
	}

	logger.Info("PostgreSQL Connection success...")

	db := &PostgresDB{
		pool: pool,
	}

	logger.Info("Pining to PostgreSQL...")
	err = db.Ping(ctx)
	if err != nil {
		logger.Error("PostgreSQL Ping Error: ", err)
		return nil, err
	}

	logger.Info("PostgreSQL Ping successful")
	return db, nil
}

func (db *PostgresDB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *PostgresDB) Close() {
	logger.Info("Closing PostgreSQL connection...")
	db.pool.Close()
	logger.Info("PostgreSQL connection closed successfully.")
}

func (db *PostgresDB) GetPool() *pgxpool.Pool {
	return db.pool
}
