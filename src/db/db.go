package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v4/pgxpool"
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
	var err error
	db.Pool, err = pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password,
		host, port, name))
	if err != nil {
		slog.Debug("cant connet to db", "error", err)
		return err
	}
	slog.Info("connection success")
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
