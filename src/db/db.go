package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type DB interface {
	Connect()
	Close()
	GetPool() *pgxpool.Pool
}

type PoolHolder struct {
	Pool *pgxpool.Pool
}

var PHolder PoolHolder

func (db *PoolHolder) Connect() error {
	user := "user_db"
	password := "1234"
	name := "songs_db"
	port := "5432"
	host := "localhost"

	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found")
	} else {
		user = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		name = os.Getenv("DB_NAME")
		port = os.Getenv("DB_PORT")
		host = os.Getenv("DB_HOST")
	}

	var err error
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
	db.Pool, err = pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		slog.Error("cant connet to db", "error", err)
		return err
	}
	slog.Info("connection success")

	if err := migrateDB(connectionString); err != nil {
		slog.Error("cant migrate", "error", err)
		return err
	}

	return nil
}
func (db *PoolHolder) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
	slog.Info("Database connection close")
}
func (db *PoolHolder) GetPool() *pgxpool.Pool {
	slog.Info("get pools")
	return db.Pool
}

func migrateDB(connectionString string) error {
	sqlDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error("connection failed", "error", err)
		return err
	}
	defer sqlDB.Close()

	m, err := migrate.New(
		"file://migrations",
		connectionString)
	if err != nil {
		slog.Error("failed to create migration", "error", err)
		return err
	}
	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		slog.Error("drop failed", "error", err)
		return err
	}
	// Применяем миграции
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		slog.Error("migration failed", "error", err)
		return err
	}

	slog.Info("migrations successful")
	return nil

}
