package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/SV1Stail/posts_and_comments/datasongs"
	"github.com/SV1Stail/posts_and_comments/db"
	mutatesong "github.com/SV1Stail/posts_and_comments/mutateSong"
	"github.com/SV1Stail/posts_and_comments/textsongs"
)

// @title Songs API
// @version 1.0
// @description API для управления песнями (CRUD операции, получение текста)
// @host localhost:8080
// @BasePath /api
func main() {
	logLevel := &slog.LevelVar{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	if os.Getenv("LOGGER") == "debug" {
		logLevel.Set(slog.LevelDebug) // установка уровня логов
	} else {
		logLevel.Set(slog.LevelInfo)
	}

	slog.SetDefault(logger)

	if err := db.PHolder.Connect(); err != nil {
		slog.Debug("connection failed")
		os.Exit(1)
	}
	defer db.PHolder.Close()
	rootMux := http.NewServeMux()

	rootMux.HandleFunc("/api/delete_song", mutatesong.Delete)
	rootMux.HandleFunc("/api/change_data", mutatesong.Change)
	rootMux.HandleFunc("/api/add_new", mutatesong.Add)

	rootMux.HandleFunc("/api/get_song_text", textsongs.Get)
	rootMux.HandleFunc("/api/get_songs", datasongs.Get)

	fmt.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", rootMux); err != nil {
		log.Fatalf("unable to up server: %v", err)
	}
}
