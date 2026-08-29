package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DB_SCHEMA = "savara_umosan"

func NewDatabase(dsn string) (*gorm.DB, error) {
	pgxConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	sqlDB := stdlib.OpenDB(
		*pgxConfig,

		// Runs once for EVERY new physical DB connection.
		stdlib.OptionAfterConnect(
			func(ctx context.Context, conn *pgx.Conn) error {
				_, err := conn.Exec(
					ctx,
					`SET search_path TO savara_umosan, public`,
				)

				if err != nil {
					return fmt.Errorf(
						"failed to set search_path: %w",
						err,
					)
				}

				return nil
			},
		),
	)

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			Conn: sqlDB,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),

			SkipDefaultTransaction: true,

			PrepareStmt: true,
		},
	)

	if err != nil {
		_ = sqlDB.Close()

		return nil, fmt.Errorf(
			"failed to connect to database: %w",
			err,
		)
	}

	// Check schema exists.
	var schemaExists bool

	err = db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM pg_namespace
			WHERE nspname = ?
		)
	`, DB_SCHEMA).Scan(&schemaExists).Error

	if err != nil {
		_ = sqlDB.Close()

		return nil, fmt.Errorf(
			"failed to verify database schema: %w",
			err,
		)
	}

	if !schemaExists {
		_ = sqlDB.Close()

		return nil, fmt.Errorf(
			"database schema %q does not exist",
			DB_SCHEMA,
		)
	}

	// Verify search path.
	var searchPath string

	if err := db.Raw(`
		SHOW search_path
	`).Scan(&searchPath).Error; err != nil {
		_ = sqlDB.Close()

		return nil, fmt.Errorf(
			"failed to read search_path: %w",
			err,
		)
	}

	var currentSchema string

	if err := db.Raw(`
		SELECT current_schema()
	`).Scan(&currentSchema).Error; err != nil {
		_ = sqlDB.Close()

		return nil, fmt.Errorf(
			"failed to read current schema: %w",
			err,
		)
	}

	// Useful while debugging the proxy/pool connection.
	var currentUser string

	if err := db.Raw(`
		SELECT current_user
	`).Scan(&currentUser).Error; err != nil {
		_ = sqlDB.Close()

		return nil, fmt.Errorf(
			"failed to read current user: %w",
			err,
		)
	}

	log.Printf(
		"[DATABASE] user=%s schema=%s search_path=%s",
		currentUser,
		currentSchema,
		searchPath,
	)

	if currentSchema != DB_SCHEMA {
		_ = sqlDB.Close()

		return nil, fmt.Errorf(
			"wrong database schema: expected %q but PostgreSQL selected %q (search_path=%q)",
			DB_SCHEMA,
			currentSchema,
			searchPath,
		)
	}

	log.Printf(
		"[DATABASE] connected successfully to schema=%s",
		DB_SCHEMA,
	)

	return db, nil
}
