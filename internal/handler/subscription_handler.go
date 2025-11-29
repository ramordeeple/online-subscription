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
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	uc *usecase.SubscriptionUseCase
}

func NewSubscriptionHandler(uc *usecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{uc: uc}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := parseCreateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := buildSubscriptionModel(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
		return
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

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseDate(str string) (int, int, error) {
	str = strings.TrimSpace(str)
	t, err := time.Parse("01-2006", str)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid date format, expected MM-YYYY")
	}
	return int(t.Month()), t.Year(), nil
}

func validateSubscriptionDates(startMonth, startYear int, endMonth, endYear *int) error {
	if endMonth != nil && endYear != nil {
		start := startYear*12 + startMonth
		end := *endYear*12 + *endMonth
		if end < start {
			return fmt.Errorf("end_date cannot be earlier than start_date")
		}
	}
	return nil
}

func parseCreateRequest(r *http.Request) (*dto.CreateSubscriptionRequest, error) {
	var req dto.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if req.UserID == nil || *req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if _, err := uuid.Parse(*req.UserID); err != nil {
		return nil, fmt.Errorf("user_id must be valid UUID")
	}

	if _, _, err := parseDate(req.StartDate); err != nil {
		return nil, fmt.Errorf("invalid start_date format, expected MM-YYYY")
	}

	if req.EndDate != nil && *req.EndDate != "" {
		if _, _, err := parseDate(*req.EndDate); err != nil {
			return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY")
		}
	}

	return &req, nil
}

func buildSubscriptionModel(req *dto.CreateSubscriptionRequest) (*model.Subscription, error) {
	startMonth, startYear, _ := parseDate(req.StartDate)

	var endMonth *int
	var endYear *int
	if req.EndDate != nil && *req.EndDate != "" {
		m, y, _ := parseDate(*req.EndDate)
		endMonth = &m
		endYear = &y
	}

	if err := validateSubscriptionDates(startMonth, startYear, endMonth, endYear); err != nil {
		return nil, err
	}

	return &model.Subscription{
		ID:          uuid.New().String(),
		UserID:      *req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartMonth:  startMonth,
		StartYear:   startYear,
		EndMonth:    endMonth,
		EndYear:     endYear,
	}, nil
}
