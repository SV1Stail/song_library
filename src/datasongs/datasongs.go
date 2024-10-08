package datasongs

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/SV1Stail/posts_and_comments/db"
	"github.com/SV1Stail/posts_and_comments/model"
	"github.com/jackc/pgx/v4"
)

// Get song's data with filtration and pagination
// @Summary Получение данных песен
// @Description Получение данных песен с фильтрацией по полям и пагинацией
// @Tags songs
// @Accept  json
// @Produce  json
// @Param page query int true "Номер страницы"
// @Param limit query int true "Количество песен на странице"
// @Param group query string false "Фильтрация по группе"
// @Param song query string false "Фильтрация по названию песни"
// @Param release_date query string false "Фильтрация по дате релиза (DD-MM-YYYY)"
// @Param text query string false "Фильтрация по тексту песни"
// @Param link query string false "Фильтрация по ссылке"
// @Success 200 {array} model.SongExtended "Список песен"
// @Failure 400 {string} string "Неверные параметры запроса"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/get_songs [get]
func Get(w http.ResponseWriter, r *http.Request) {
	slog.Info("request to get song info")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		slog.Error("wrong page number", "page", page, "error", err)
		http.Error(w, "wrong page number", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		slog.Error("wrong limit number", "limit", limit, "error", err)
		http.Error(w, "wrong limit number", http.StatusBadRequest)
		return
	}
	slog.Info("conver successful", "page", page, "limit", limit)

	var song model.SongExtended

	song.Band = r.URL.Query().Get("group")
	song.Name = r.URL.Query().Get("song")
	song.Link = r.URL.Query().Get("link")
	song.RDate = r.URL.Query().Get("release_date")
	song.Text = r.URL.Query().Get("text")

	slog.Info("get params successful")

	ctx := context.Background()
	conn, err := db.PHolder.GetPool().Acquire(ctx)
	if err != nil {
		slog.Error("cant make conn", "error", err)
		http.Error(w, "cant make conn", http.StatusInternalServerError)
		return
	}
	defer conn.Release()
	slog.Info("made conn from connections pool")

	ta, err := conn.Begin(ctx)
	if err != nil {
		slog.Error("cant begin transaction", "error", err)
		http.Error(w, "cant begin transaction", http.StatusInternalServerError)

		return
	}
	slog.Info("transaction in work")

	defer ta.Rollback(ctx)
	DBrequest, args := makeStringForDBRequest(&song, page, limit)
	rows, err := ta.Query(ctx, DBrequest, args...)
	if err == pgx.ErrNoRows {
		slog.Error("no rows with your params")
		http.Error(w, "no rows with your params", http.StatusInternalServerError)
		return
	} else if err != nil {
		slog.Error("error in SELECT operation", "error", err)
		http.Error(w, "cant SELECT from db", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	slog.Info("rows collected")

	var songs []model.SongExtended
	for rows.Next() {
		if err := rows.Scan(&song.Band, &song.Name, &song.RDate, &song.Text, &song.Link); err != nil {
			slog.Error("cant Scan row from db", "error", err)
			http.Error(w, "cant Scan row from db", http.StatusInternalServerError)
			return
		}
		songs = append(songs, song)
		slog.Debug("add song", "Band", song.Band, "song", song.Name)
	}
	if err := rows.Err(); err != nil {
		slog.Error("smt go wrong, rows have err", "error", err)
		http.Error(w, "cant Scan row from db", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		slog.Error("cant encode JSON", "error", err)
		http.Error(w, "cant encode JSON", http.StatusInternalServerError)
		return
	}
	slog.Info("response sent")

}

// mkae string for request to DB
func makeStringForDBRequest(song *model.SongExtended, page, limit int) (string, []interface{}) {
	DBrequest := "SELECT Band,song,release_date,text,link FROM songs_table WHERE 1=1"
	args := []interface{}{}
	argID := 1
	if song.Band != "" {
		DBrequest += "AND Band=$" + strconv.Itoa(argID)
		args = append(args, song.Band)
		argID++
	}
	if song.Name != "" {
		DBrequest += " AND song=$" + strconv.Itoa(argID)
		args = append(args, song.Name)
		argID++
	}
	if song.RDate != "" {
		DBrequest += " AND release_date=$" + strconv.Itoa(argID)
		args = append(args, song.RDate)
		argID++
	}
	if song.Text != "" {
		DBrequest += " AND text=$%" + strconv.Itoa(argID)
		args = append(args, song.Text)
		argID++
	}
	if song.Link != "" {
		DBrequest += " AND link=$" + strconv.Itoa(argID)
		args = append(args, song.Link)
		argID++
	}

	offset := (page - 1) * limit
	DBrequest += " LIMIT $" + strconv.Itoa(argID) + " OFFSET $" + strconv.Itoa(argID+1)
	args = append(args, limit, offset)

	return DBrequest, args

}
