package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/SV1Stail/posts_and_comments/db"
	mutatesong "github.com/SV1Stail/posts_and_comments/mutateSong"
)

func main() {
	logLevel := &slog.LevelVar{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	logLevel.Set(slog.LevelInfo) // установка уровня логов

	slog.SetDefault(logger)

	if err := db.PHolder.Connect(); err != nil {
		slog.Debug("connection failed")
		os.Exit(1)
	}
	defer db.PHolder.Close()
	rootMux := http.NewServeMux()
	songsMux := http.NewServeMux()

	rootMux.Handle("/api", http.StripPrefix("/api", songsMux))
	songsMux.HandleFunc("/delete_song", mutatesong.Delete)
	songsMux.HandleFunc("/change_data", mutatesong.Change)
	songsMux.HandleFunc("/add_new", mutatesong.Add)

	// rootMux.HandleFunc("/api/get_song_text", GetSongText)
	// rootMux.HandleFunc("/api/get_songs", GetSongs)

	fmt.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", rootMux); err != nil {
		log.Fatalf("unable to up server: %v", err)
	}
}
