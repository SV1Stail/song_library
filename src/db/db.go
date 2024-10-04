package db

import (
	"context"
	"fmt"
	"log"

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

func (db *PoolHolder) Connect() {
	user := "user_db"
	password := "1234"
	name := "songs_db"
	port := "5432"
	host := "localhost"
	var err error
	db.Pool, err = pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password,
		host, port, name))
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	log.Println("Connected to database successfully")
}
func (db *PoolHolder) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
	log.Println("Database connection close")
}
func (db *PoolHolder) GetPool() *pgxpool.Pool {
	return db.Pool
}
