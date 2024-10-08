package textsongs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/SV1Stail/posts_and_comments/db"
	"github.com/SV1Stail/posts_and_comments/model"
)

// get song's text and verse by verse pagination
// @Summary Получение текста песни
// @Description Получение текста песни с пагинацией по куплетам. Фильтрация по названию песни и группы, страницам и лимитом для постраничного отображения куплетов.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param page query int true "Номер страницы"
// @Param limit query int true "Количество песен на странице"
// @Param group query string false "Фильтрация по группе"
// @Param song query string false "Фильтрация по названию песни"
// @Success 200 {array} string "Список куплетов"
// @Failure 400 {string} string "Неверные параметры запроса)"
// @Failure 404 {string} string "Количество куплетов превышает доступные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/get_song_text [get]
func Get(w http.ResponseWriter, r *http.Request) {
	slog.Info("request to get song text")
	var song model.SongExtended
	var couplets []string
	var start, end, coupletsLen int

	song.Name = r.URL.Query().Get("song")
	song.Band = r.URL.Query().Get("group")

	slog.Info("decode successful")

	ctx := context.Background()
	if err := song.GetSongFromDB(ctx, db.PHolder.GetPool()); err != nil {
		slog.Error("cant get song from db", "error", err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	slog.Info("got song from db")

	if song.Text == "" {
		slog.Warn("song does not have text")
		http.Error(w, "song does not have text", http.StatusBadRequest)
		return
	}
	slog.Info("song have text")

	couplets = strings.Split(song.Text, "\n\n")
	slog.Info("couplets split")

	coupletsLen = len(couplets)
	slog.Info("Retrieved couplets length", "length", coupletsLen)
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
	slog.Info("Valid page and limit parameters", "page", page, "limit", limit)

	start = (page - 1) * limit
	end = start + limit
	slog.Info("start and end indices", "start", start, "end", end)

	if start >= coupletsLen {
		slog.Error("song does not have to much couplets", "requested_start", start, "page", page, "total_couplets", coupletsLen)
		http.Error(w, fmt.Sprintf("song does not have %d couplets on page %d", start, page), http.StatusNotFound)
		return
	}
	if end > coupletsLen {
		end = coupletsLen
	}
	necesaryCouplets := couplets[start:end]
	resp, err := json.Marshal(necesaryCouplets)
	if err != nil {
		slog.Error("cnt marshal response ", "error", err)

		http.Error(w, "cant marshal response", http.StatusInternalServerError)
		return
	}
	slog.Info("ready to send")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	slog.Info("response sent")
}
