package mutatesong

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/SV1Stail/posts_and_comments/db"
	"github.com/SV1Stail/posts_and_comments/model"
)

// delete song from DB
// @Summary Удаление данных песни
// @Description Удаление из БД строки с песней по названию и группе
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group query string false "Фильтрация по группе"
// @Param song query string false "Фильтрация по названию песни"
// @Success 200 {string} string "Delete was successful"
// @Failure 400 {string} string "Неверные параметры запроса"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/delete_song [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	slog.Info("request Delete song")

	var song model.SongExtended
	song.Band = r.URL.Query().Get("group")
	song.Name = r.URL.Query().Get("song")

	if song.Band == "" || song.Name == "" {
		slog.Error("song name or group is empty!", "song", song.Name, "group", song.Band)
		http.Error(w, "check song and group fields", http.StatusBadRequest)
		return
	}

	slog.Info("data collected from request")
	pool := db.PHolder.GetPool()
	slog.Debug("get pool of connections successful")
	ctx := context.Background()
	slog.Debug("context made")
	if err := song.DeleteFromDB(ctx, pool); err != nil {
		slog.Error("DELETE FROM failed", "error", err)

		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	slog.Info("DELETE successful")

	w.WriteHeader(http.StatusOK)
	slog.Info("header 200")

	w.Write([]byte("Delete was successful"))
	slog.Info("response made")

}

// @Summary Изменение данных песни
// @Description Изменение данных песни (дата релиза, текст, ссылка). Поиск по группе и названию.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group query string false "Фильтрация по группе"
// @Param song query string false "Фильтрация по названию песни"
// @Param release_date query string false "Дата релиза (DD-MM-YYYY)"
// @Param text query string false "Текст песни"
// @Param link query string false "Ссылка"
// @Success 200 {string} string "Update was successful"
// @Failure 400 {string} string "Неверные параметры запроса"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/change_data [put]
func Change(w http.ResponseWriter, r *http.Request) {
	slog.Info("request Change song")

	var song model.SongExtended

	song.Band = r.URL.Query().Get("group")
	song.Name = r.URL.Query().Get("song")
	song.RDate = r.URL.Query().Get("release_date")
	song.Text = r.URL.Query().Get("text")
	song.Link = r.URL.Query().Get("link")

	if song.Band == "" || song.Name == "" {
		slog.Error("song name or group is empty!", "song", song.Name, "group", song.Band)
		http.Error(w, "check song and group fields", http.StatusBadRequest)
		return
	}

	slog.Info("data collected from request")

	if song.RDate == "" && song.Text == "" && song.Link == "" {
		slog.Warn("no need to change smt")
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}
	slog.Info("data is valid")

	pool := db.PHolder.GetPool()
	slog.Debug("get pool of connections successful")
	ctx := context.Background()
	slog.Debug("context made")
	if err := song.ChangeInDB(ctx, pool); err != nil {
		slog.Error("CHANGE data in db failed", "error", err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	slog.Info("CHANGE successful")

	w.WriteHeader(http.StatusOK)
	slog.Info("header 200")

	w.Write([]byte("Update was successful"))
	slog.Info("response made")

}

// add new song in DB
// @Summary Добавление новой песни
// @Description Добавление новой песни в базу данных с помощью внешнего API
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body model.SongExtended true "Данные новой песни"
// @Success 201 {string} string "Add successful"
// @Failure 400 {string} string "Неверные параметры запроса"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/add_new [post]
func Add(w http.ResponseWriter, r *http.Request) {
	slog.Info("request Add song")

	var song model.SongExtended
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		slog.Error("wrong request body", "error", err)
		http.Error(w, "wrong request body", http.StatusBadRequest)
		return
	}
	slog.Info("data collected from request")
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://example.com/info"
	}
	slog.Debug("adding params")
	params := url.Values{}
	params.Add("song", song.Name)
	params.Add("group", song.Band)
	slog.Info("params added", "song", song.Name, "Band", song.Band)

	resp, err := http.Get(apiURL + "?" + params.Encode())
	slog.Debug("have full response", "resp", resp)
	if err != nil {
		slog.Error("error when get request to api", "error", err)
		http.Error(w, "error when get request to api", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	slog.Info("response successful")

	if resp.StatusCode != http.StatusOK {
		slog.Error("error from api", "code", resp.StatusCode, "error", err)
		http.Error(w, "error from api", resp.StatusCode)
		return
	}
	slog.Debug("reading body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("reading response body", "error", err)

		http.Error(w, "reading response body", http.StatusInternalServerError)
		return
	}
	slog.Debug("reading body successful")

	slog.Debug("Unmarshal json")
	err = json.Unmarshal(body, &song)
	if err != nil {
		slog.Error("unmarshal was failed", "error", err)
		http.Error(w, "unmarshal was failed", http.StatusInternalServerError)
		return
	}
	slog.Debug("Unmarshal json successful")

	pool := db.PHolder.GetPool()
	slog.Debug("get pool of connections successful")
	ctx := context.Background()
	slog.Debug("context made")
	if err := song.SaveInDB(ctx, pool); err != nil {
		slog.Error("SAVE failed", "error", err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	slog.Info("SAVE successful")

	w.WriteHeader(http.StatusCreated)
	slog.Info("header 201")

	w.Write([]byte("Add successful"))
	slog.Info("response made")

}
