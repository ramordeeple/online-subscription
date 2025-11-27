package handler

import (
	"encoding/json"
	"net/http"
	"online-subscription/internal/model"
	"online-subscription/internal/repository"
	"online-subscription/internal/usecase"
	"time"
)

type SubscriptionHandler struct {
	uc *usecase.SubscriptionUseCase
}

func NewSubscriptionHandler(uc *usecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{uc: uc}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var s model.Subscription
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.uc.Create(r.Context(), &s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, s)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	f := repository.SubscriptionFilter{
		UserID:      ptrString(r.URL.Query().Get("user_id")),
		ServiceName: ptrString(r.URL.Query().Get("service_name")),
	}

	subs, err := h.uc.List(r.Context(), &f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, subs)
}

func (h *SubscriptionHandler) GetById(w http.ResponseWriter, r *http.Request, id string) {
	s, err := h.uc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if s == nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
	}

	writeJSON(w, http.StatusOK, s)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) Summary(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	fromMonth, fromYear, err := parseDate(from)
	if err != nil {
		http.Error(w, "invalid to date", http.StatusBadRequest)
		return
	}

	toMonth, toYear, err := parseDate(to)
	if err != nil {
		http.Error(w, "invalid to date", http.StatusBadRequest)
		return
	}

	f := repository.SummaryFilter{
		FromMonth:   fromMonth,
		FromYear:    fromYear,
		ToMonth:     toMonth,
		ToYear:      toYear,
		UserID:      ptrString(r.URL.Query().Get("user_id")),
		ServiceName: ptrString(r.URL.Query().Get("service_name")),
	}

	sum, err := h.uc.Sum(r.Context(), &f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"total": sum})
}

func parseDate(str string) (int, int, error) {
	t, err := time.Parse("01-2006", str)
	if err != nil {
		return 0, 0, err
	}
	return int(t.Month()), t.Year(), nil
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
