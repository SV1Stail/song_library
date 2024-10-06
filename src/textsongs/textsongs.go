package textsongs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/SV1Stail/posts_and_comments/db"
	"github.com/SV1Stail/posts_and_comments/model"
)

// get song's text and verse by verse pagination
// params: ?page, limit
func Get(w http.ResponseWriter, r *http.Request) {
	var song model.SongExtended
	var couplets []string
	var wg sync.WaitGroup
	var start, end, coupletsLen int
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
			http.Error(w, fmt.Sprintf("cant decode json %v", err), http.StatusBadRequest)
			return
		}
		ctx := context.Background()
		if err := song.GetSongFromDB(ctx, db.PHolder.GetPool()); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		if song.Text == "" {
			http.Error(w, "song does not have text", http.StatusBadRequest)
			return
		} else {
			couplets = strings.Split(song.Text, "\n\n")
		}
		coupletsLen = len(couplets)
	}()

	go func() {
		defer wg.Done()
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			http.Error(w, fmt.Sprintf("wrong page number %v", err), http.StatusBadRequest)
			return
		}
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit < 1 {
			http.Error(w, fmt.Sprintf("wrong limit number %v", err), http.StatusBadRequest)
			return
		}
		start = (page - 1) * limit
		end = start + limit
	}()

	wg.Wait()

	if start > coupletsLen {
		http.Error(w, fmt.Sprintf("song does not have %d couplets", start), http.StatusNotFound)
		return
	}
	if end > coupletsLen {
		end = coupletsLen
	}
	nessesaryCouplets := couplets[start:end]
	resp, err := json.Marshal(nessesaryCouplets)
	if err != nil {
		http.Error(w, fmt.Sprintf("cnt marshal response | %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
