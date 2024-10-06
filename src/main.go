package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SV1Stail/posts_and_comments/db"
	mutatesong "github.com/SV1Stail/posts_and_comments/mutateSong"
)

func main() {

	db.PHolder.Connect()
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
