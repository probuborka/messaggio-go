package handlerhttp

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/probuborka/messaggio/internal/domain"
	"github.com/probuborka/messaggio/internal/service"
	"github.com/probuborka/messaggio/pkg/logger"
)

type Handler struct {
	message service.Message
}

func New(services *service.Services) *Handler {
	return &Handler{
		message: services.Message,
	}
}

func (h Handler) Init() *chi.Mux {
	r := chi.NewRouter()

	// Создать сообщение
	r.Post("/message", h.postMessage)

	// Получить статистику
	r.Get("/statistics", h.getStatistics)

	return r
}

func (h Handler) postMessage(w http.ResponseWriter, r *http.Request) {
	var message domain.Message
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error(err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error(err)
		return
	}

	if err := h.message.Create(r.Context(), message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h Handler) getStatistics(w http.ResponseWriter, r *http.Request) {
	message, err := h.message.Statistics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error(err)
		return
	}

	resp, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
