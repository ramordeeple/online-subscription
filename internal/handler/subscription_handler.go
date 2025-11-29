package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"online-subscription/internal/handler/dto"
	"online-subscription/internal/logger"
	"online-subscription/internal/model"
	"online-subscription/internal/usecase"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
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
	var req dto.CreateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startMonth, startYear, err := parseDate(req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format, expected YYYY-MM", http.StatusBadRequest)
		return
	}

	var userID string
	if req.UserID != nil && *req.UserID != "" {
		userID = *req.UserID
	} else {
		userID = uuid.New().String()
	}

	sub := &model.Subscription{
		ID:          uuid.New().String(),
		UserID:      userID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartMonth:  startMonth,
		StartYear:   startYear,
	}

	if err := h.uc.Create(r.Context(), sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Subscription created",
		zap.String("id", sub.ID),
		zap.String("service", sub.ServiceName),
		zap.String("user_id", sub.UserID),
	)

	writeJSON(w, http.StatusCreated, sub)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	f := model.SubscriptionFilter{
		UserID:      ptrString(r.URL.Query().Get("user_id")),
		ServiceName: ptrString(r.URL.Query().Get("service_name")),
	}

	subs, err := h.uc.List(r.Context(), &f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Subscriptions listed", zap.Int("count", len(subs)))
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

	logger.Info("Subscription retrieved", zap.String("id", s.ID))
	writeJSON(w, http.StatusOK, s)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Subscription deleted", zap.String("id", id))
	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) Summary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	from := r.URL.Query().Get("from")

	fromMonth, fromYear, err := parseDate(from)
	if err != nil {
		http.Error(w, "invalid from date", http.StatusBadRequest)
		return
	}

	f := model.SummaryFilter{
		FromMonth:   fromMonth,
		FromYear:    fromYear,
		UserID:      ptrString(r.URL.Query().Get("user_id")),
		ServiceName: ptrString(r.URL.Query().Get("service_name")),
	}

	sum, err := h.uc.Sum(r.Context(), &f)
	if err != nil {
		logger.Error("Failed to calculate summary", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Summary calculated", zap.Int("sum", sum))
	writeJSON(w, http.StatusOK, map[string]int{"total": sum})
}

func parseDate(str string) (int, int, error) {
	str = strings.TrimSpace(str)
	t, err := time.Parse("01-2006", str)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid date format, expected MM-YYYY-MM")
	}
	return int(t.Month()), t.Year(), nil
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
