package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"online-subscription/internal/handler/helpers"
	"online-subscription/internal/handler/mapper"
	"online-subscription/internal/handler/parser"
	"online-subscription/internal/logger"
	"online-subscription/internal/model"
	"online-subscription/internal/usecase"
	"strings"

	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	uc *usecase.SubscriptionUseCase
}

func NewSubscriptionHandler(uc *usecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{uc: uc}
}

// Create godoc
// @Summary Create subscription
// @Description Creates a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} model.Subscription
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := parser.ParseCreateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := mapper.BuildSubscriptionModel(req)
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

// List godoc
// @Summary List subscriptions
// @Description Returns list of subscriptions with optional filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Success 200 {array} model.Subscription
// @Failure 500 {string} string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	f := model.SubscriptionFilter{
		UserID:      helpers.PtrString(r.URL.Query().Get("user_id")),
		ServiceName: helpers.PtrString(r.URL.Query().Get("service_name")),
	}

	subs, err := h.uc.List(r.Context(), &f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Subscriptions listed", zap.Int("count", len(subs)))
	writeJSON(w, http.StatusOK, subs)
}

// GetById godoc
// @Summary Get subscription by ID
// @Description Returns a subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [get]
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

// Delete godoc
// @Summary Delete subscription
// @Description Deletes subscription by ID
// @Tags subscriptions
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 500 {string} string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.uc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Subscription deleted", zap.String("id", id))
	w.WriteHeader(http.StatusNoContent)
}

// Summary godoc
// @Summary Summary of payments
// @Description Calculates total cost of subscriptions for date range
// @Tags subscriptions
// @Produce json
// @Param from query string true "Start date MM-YYYY"
// @Param to query string false "End date MM-YYYY"
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Success 200 {object} map[string]int
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/summary [get]
func (h *SubscriptionHandler) Summary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	fromMonth, fromYear, err := helpers.ParseDate(from)
	if err != nil {
		http.Error(w, "invalid from date", http.StatusBadRequest)
		return
	}

	var toMonth, toYear *int
	if strings.TrimSpace(to) != "" {
		m, y, err := helpers.ParseDate(to)
		if err != nil {
			http.Error(w, "invalid to date", http.StatusBadRequest)
			return
		}
		if y*12+m < fromYear*12+fromMonth {
			http.Error(w, "`to` date cannot be earlier than `from` date", http.StatusBadRequest)
			return
		}
		toMonth = &m
		toYear = &y
	}

	f := model.SummaryFilter{
		FromMonth:   fromMonth,
		FromYear:    fromYear,
		ToMonth:     toMonth,
		ToYear:      toYear,
		UserID:      helpers.PtrString(r.URL.Query().Get("user_id")),
		ServiceName: helpers.PtrString(r.URL.Query().Get("service_name")),
	}

	sum, err := h.uc.Sum(r.Context(), &f)
	if err != nil {
		logger.Error("Failed to calculate summary", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Summary calculated",
		zap.Int("sum", sum),
		zap.String("user_id", helpers.SafeString(f.UserID)),
		zap.String("service_name", helpers.SafeString(f.ServiceName)),
		zap.String("from", fmt.Sprintf("%02d-%d", fromMonth, fromYear)),
		zap.String("to", func() string {
			if toMonth != nil && toYear != nil {
				return fmt.Sprintf("%02d-%d", *toMonth, *toYear)
			}
			return ""
		}()),
	)

	writeJSON(w, http.StatusOK, map[string]int{"total": sum})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
