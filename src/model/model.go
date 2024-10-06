package model

import (
	"context"
	"fmt"

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
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("cant make conn %v", err)
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, "SELECT release_date,text,link FROM songs_table WHERE group=$1 AND name=$2", song.Group, song.Name).Scan(&song.RDate, &song.Text, &song.Link)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("no rows with this group: %s and song %s", song.Group, song.Name)
	} else if err != nil {
		return fmt.Errorf("cant SELECT from db %v", err)
	}

	return nil
}

func (song *SongExtended) SaveInDB(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("cant make conn: %v", err)
	}
	defer conn.Release()
	ta, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cant begin transaction: %v", err)
	}

	defer ta.Rollback(ctx)
	_, err = ta.Exec(ctx, "INSERT INTO songs_table (group, song, release_date, text, link) VALUES $1,$2,$3,$4,$5",
		song.Group, song.Name, song.RDate, song.Text, song.Link)
	if err != nil {
		return fmt.Errorf("cant INSERT INTO songs_table: %v", err)
	}
	err = ta.Commit(ctx)
	if err != nil {
		return fmt.Errorf("cant commit tranaction: %v", err)
	}
	return nil
}
func (song *SongExtended) DeleteFromDB(ctx context.Context, pool *pgxpool.Pool) error {
	coon, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("cant make conn: %v", err)
	}
	defer coon.Release()
	ta, err := coon.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cant begin transaction: %v", err)
	}
	defer ta.Rollback(ctx)
	_, err = ta.Exec(ctx, "DELETE FROM songs_table WHERE group=$1 AND song=$2", song.Group, song.Name)
	if err != nil {
		return fmt.Errorf("cant DELETE FROM songs_table: %v", err)
	}
	err = ta.Commit(ctx)
	if err != nil {
		return fmt.Errorf("cant commit tranaction: %v", err)
	}
	return nil
}
func (song *SongExtended) ChangeInDB(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("cant make conn: %v", err)

	}
	defer conn.Release()
	ta, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cant begin transaction: %v", err)
	}
	defer ta.Rollback(ctx)

	query := "UPDATE songs_table SET "
	args := []interface{}{}
	argID := 1
	if song.RDate != "" {
		query += fmt.Sprintf("release_date=$%d, ", argID)
		args = append(args, song.RDate)
		argID++
	}
	if song.Text != "" {
		query += fmt.Sprintf("text=$%d, ", argID)
		args = append(args, song.Text)
		argID++
	}
	if song.Link != "" {
		query += fmt.Sprintf("link=$%d, ", argID)
		args = append(args, song.Link)
		argID++
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf("WHERE group=$%d AND song=$%d", argID, argID+1)
	args = append(args, song.Group, song.Name)

	_, err = ta.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cant update database: %v", err)
	}
	err = ta.Commit(ctx)
	if err != nil {
		return fmt.Errorf("cant commit tranaction: %v", err)
	}
	return nil
}
