package model

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Song struct {
	Group string `json:"group"`
	Name  string `json:"song"`
}
type SongExtended struct {
	Song
	RDate string `json:"release_date,omitempty"`
	Text  string `json:"text,omitempty"`
	Link  string `json:"link,omitempty"`
}

func (song *SongExtended) GetSongFromDB(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("start GetSongFromDB")
	conn, err := pool.Acquire(ctx)
	if err != nil {
		slog.Error("cant make conn", "error", err)
		return fmt.Errorf("cant make conn %v", err)
	}
	defer conn.Release()
	slog.Info("made conn from connections pool")

	err = conn.QueryRow(ctx, "SELECT release_date,text,link FROM songs_table WHERE group=$1 AND name=$2", song.Group, song.Name).Scan(&song.RDate, &song.Text, &song.Link)
	if err == pgx.ErrNoRows {
		slog.Error("no rows with tihs group nad song", "group", song.Group, "song", song.Name)
		return fmt.Errorf("no rows with this group: %s and song %s", song.Group, song.Name)
	} else if err != nil {
		slog.Error("error in SELECT operation", "error", err)
		return fmt.Errorf("cant SELECT from db")
	}
	slog.Info("SELECT successful")

	return nil
}

func (song *SongExtended) SaveInDB(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("start SaveInDB")

	conn, err := pool.Acquire(ctx)
	if err != nil {
		slog.Error("cant make conn", "error", err)
		return fmt.Errorf("cant make conn")
	}
	defer conn.Release()
	slog.Info("made conn from connections pool")

	ta, err := conn.Begin(ctx)
	if err != nil {
		slog.Error("cant begin ta", "error", err)
		return fmt.Errorf("cant begin transaction")
	}
	slog.Info("transaction in work")

	defer ta.Rollback(ctx)
	_, err = ta.Exec(ctx, "INSERT INTO songs_table (group, song, release_date, text, link) VALUES $1,$2,$3,$4,$5",
		song.Group, song.Name, song.RDate, song.Text, song.Link)
	if err != nil {
		slog.Error("cant INSERT in db", "error", err)
		return fmt.Errorf("cant INSERT INTO songs_table")
	}
	slog.Info("INSERT INTO done")
	err = ta.Commit(ctx)
	if err != nil {
		slog.Error("cant commit", "error", err)
		return fmt.Errorf("cant commit tranaction")
	}
	slog.Info("commit made")

	return nil
}
func (song *SongExtended) DeleteFromDB(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("start DeleteFromDB")
	coon, err := pool.Acquire(ctx)
	if err != nil {
		slog.Error("cant make conn", "error", err)

		return fmt.Errorf("cant make conn")
	}
	defer coon.Release()
	slog.Info("made conn from connections pool")

	ta, err := coon.Begin(ctx)
	if err != nil {
		slog.Error("cant begin ta", "error", err)

		return fmt.Errorf("cant begin transaction")
	}
	slog.Info("transaction in work")

	defer ta.Rollback(ctx)
	_, err = ta.Exec(ctx, "DELETE FROM songs_table WHERE group=$1 AND song=$2", song.Group, song.Name)
	if err != nil {
		slog.Error("cant DELETE FROM db", "error", err)

		return fmt.Errorf("cant DELETE FROM songs_table")
	}
	slog.Info("DELETE FROM songs_table successful")

	err = ta.Commit(ctx)
	if err != nil {
		slog.Error("cant commit tranaction", "error", err)
		return fmt.Errorf("cant commit tranaction")
	}
	slog.Info("commit tranaction successful")

	return nil
}
func (song *SongExtended) ChangeInDB(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("start ChangeInDB")

	conn, err := pool.Acquire(ctx)
	if err != nil {
		slog.Error("cant make conn", "error", err)

		return fmt.Errorf("cant make conn")
	}
	defer conn.Release()
	slog.Info("made conn from connections pool")

	ta, err := conn.Begin(ctx)
	if err != nil {
		slog.Error("cant begin ta", "error", err)
		return fmt.Errorf("cant begin transaction")
	}
	slog.Info("transaction in work")

	defer ta.Rollback(ctx)

	query := "UPDATE songs_table SET "
	args := []interface{}{}
	argID := 1
	if song.RDate != "" {
		slog.Debug("need release_date")
		query += "release_date=$" + strconv.Itoa(argID) + ", "
		args = append(args, song.RDate)
		argID++
	}
	if song.Text != "" {
		slog.Debug("need text")
		query += "text=$" + strconv.Itoa(argID) + ", "
		args = append(args, song.Text)
		argID++
	}
	if song.Link != "" {
		slog.Debug("need link")
		query += "link=$" + strconv.Itoa(argID) + ", "
		args = append(args, song.Link)
		argID++
	}
	slog.Info("get what we have to update")

	query = query[:len(query)-2]
	query += fmt.Sprintf("WHERE group=$%d AND song=$%d", argID, argID+1)
	args = append(args, song.Group, song.Name)

	slog.Info("STRING FRO REQUEST READY")
	_, err = ta.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("cant update database", "error", err)
		return fmt.Errorf("cant update database")
	}
	slog.Info("request for update")

	err = ta.Commit(ctx)
	if err != nil {
		slog.Error("cant commit tranaction", "error", err)
		return fmt.Errorf("cant commit tranaction")
	}
	slog.Info("commit tranaction successful")
	return nil
}
