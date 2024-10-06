package mutatesong

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/SV1Stail/posts_and_comments/db"
	"github.com/SV1Stail/posts_and_comments/model"
)

// delete song from DB
func Delete(w http.ResponseWriter, r *http.Request) {
	var song model.SongExtended
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, fmt.Sprintf("wrong request body: %v", err), http.StatusBadRequest)
		return
	}
	pool := db.PHolder.GetPool()
	ctx := context.Background()
	if err := song.DeleteFromDB(ctx, pool); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete was successful"))

}

// change song's data in DB
func Change(w http.ResponseWriter, r *http.Request) {

	var song model.SongExtended
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, fmt.Sprintf("wrong request body: %v", err), http.StatusBadRequest)
		return
	}

	if song.RDate == "" && song.Text == "" && song.Link == "" {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	pool := db.PHolder.GetPool()
	ctx := context.Background()
	if err := song.ChangeInDB(ctx, pool); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Update was successful"))
}

// add new song in DB
func Add(w http.ResponseWriter, r *http.Request) {
	var song model.SongExtended
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, fmt.Sprintf("wrong request body: %v", err), http.StatusBadRequest)
		return
	}

	apiURL := "http://example.com/info"
	params := url.Values{}
	params.Add("song", song.Name)
	params.Add("group", song.Group)

	resp, err := http.Get(apiURL + "?" + params.Encode())
	if err != nil {
		http.Error(w, fmt.Sprintf("error when get request to api: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("error from api %d", resp.StatusCode), resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("reading response body %v", err), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &song)
	if err != nil {
		http.Error(w, fmt.Sprintf("unmarshal was failed %v", err), http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	pool := db.PHolder.GetPool()
	if err := song.SaveInDB(ctx, pool); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Add successful"))
}
