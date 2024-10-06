package mutatesong

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SV1Stail/posts_and_comments/db"
)

type Song struct {
	Group string `json:"group"`
	Name  string `json:"song"`
}
type SongExtended struct {
	Song
	RDate time.Time `json:"release_date,omitempty"`
	Text  string    `json:"text,omitempty"`
	Link  string    `json:"link,omitempty"`
}

// delete song from DB
func Delete(w http.ResponseWriter, r *http.Request) {
	var song Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, fmt.Sprintf("wrong request body: %v", err), http.StatusBadRequest)
		return
	}
	pool := db.PHolder.GetPool()
	ctx := context.Background()
	coon, err := pool.Acquire(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("cant make conn: %v", err), http.StatusInternalServerError)
		return
	}
	defer coon.Release()
	ta, err := coon.Begin(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("cant begin transaction: %v", err), http.StatusInternalServerError)
		return
	}
	defer ta.Rollback(ctx)
	_, err = ta.Exec(ctx, "DELETE FROM songs_table WHERE group=$1 AND song=$2", song.Group, song.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("cant DELETE FROM songs_table: %v", err), http.StatusInternalServerError)
		return
	}
	err = ta.Commit(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("cant commit tranaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete was successful"))

}

// change song's data in DB
func Change(w http.ResponseWriter, r *http.Request) {

	var song SongExtended
	var err error
	err = json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, fmt.Sprintf("wrong request body: %v", err), http.StatusBadRequest)
		return
	}

	if song.RDate.IsZero() && song.Text == "" && song.Link == "" {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	pool := db.PHolder.GetPool()
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("cant make conn: %v", err), http.StatusInternalServerError)

		return
	}
	defer conn.Release()
	ta, err := conn.Begin(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("cant begin transaction: %v", err), http.StatusInternalServerError)

		return
	}
	defer ta.Rollback(ctx)

	query := "UPDATE songs_table SET "
	args := []interface{}{}
	argID := 1
	if !song.RDate.IsZero() {
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
		http.Error(w, fmt.Sprintf("cant update database: %v", err), http.StatusInternalServerError)
		return
	}
	err = ta.Commit(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("cant commit tranaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Update was successful"))
}

// add new song in DB
func Add(w http.ResponseWriter, r *http.Request) {
	var song Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, fmt.Sprintf("wrong request body: %v", err), http.StatusBadRequest)
		return
	}

}
