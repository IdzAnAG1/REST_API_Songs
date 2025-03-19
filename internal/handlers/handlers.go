package handlers

import (
	"REST_API_Songs/internal/service"
	"REST_API_Songs/internal/structure"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Handling GET /library request")

	filters := map[string]string{
		"group_title":  r.URL.Query().Get("group_title"),
		"song_title":   r.URL.Query().Get("song_title"),
		"release_date": r.URL.Query().Get("release_date"),
		"song_text":    r.URL.Query().Get("song_text"),
		"link":         r.URL.Query().Get("link"),
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	songs, err := h.service.GetSongs(r.Context(), filters, page, limit)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch songs"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		logrus.Errorf("Failed to encode response: %v", err)
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// GetSongText Получает текст песни с пагинацией по куплетам .
// @Summary Получить текст песни с пагинацией по куплетам
// @Description Получение текста песни с постраничной разбивкой по куплетам
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of verses per page" default(2)
// @Success 200 {array} string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /song/{id}/text [get]
func (h *Handler) GetSongText(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Handling GET /song/{id}/text request")
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid song ID"}`, http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 2
	}
	verses, err := h.service.GetSongText(r.Context(), id, page, limit)
	if err != nil {
		if strings.Contains(err.Error(), "song not found") {
			http.Error(w, `{"error":"Song not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Failed to fetch song text"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(verses); err != nil {
		logrus.Errorf("Failed to encode response: %v", err)
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// DeleteSong Производит удаление песни по ID.
// @Summary Удаление песни
// @Description Удаляет песню находя ее по ее уникальному ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /song/{id} [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Handling DELETE /song/{id} request")
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid song ID"}`, http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteSong(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, `{"error":"Song not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Failed to delete song"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateSong Обнавляет информацию о песни по ее ID.
// @Summary Обновление информации о песне
// @Description Обнавляет информацию о песни по ее ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body structure.Song true "Song data"
// @Success 200 {object} structure.Song
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /song/{id} [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Handling PUT /song/{id} request")
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid song ID"}`, http.StatusBadRequest)
		return
	}
	var song structure.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateSong(r.Context(), id, &song); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, `{"error":"Song not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Failed to update song"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(song); err != nil {
		logrus.Errorf("Failed to encode response: %v", err)
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// CreateSong Создает новую запись о песне.
// @Summary Создает новую песню
// @Description Добавление новой песни и получение дополнительных сведений из внешнего API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body structure.NewSong true "Song input"
// @Success 201 {object} map[string]int64
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /song [post]
func (h *Handler) CreateSong(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Handling POST /song request")
	var in structure.NewSong
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}
	if in.GroupTitle == "" || in.SongTitle == "" {
		http.Error(w, `{"error":"Group and song are required"}`, http.StatusBadRequest)
		return
	}
	id, err := h.service.CreateSong(r.Context(), &in)
	if err != nil {
		http.Error(w, `{"error":"Failed to create song"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]int64{"id": id}); err != nil {
		logrus.Errorf("Failed to encode response: %v", err)
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}
