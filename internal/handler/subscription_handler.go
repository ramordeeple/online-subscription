package handler

import (
	"encoding/json"
	"net/http"
	"online-subscription/internal/handler/dto"
	"online-subscription/internal/handler/helpers"
	"online-subscription/internal/handler/mapper"
	"online-subscription/internal/handler/parser"
	"online-subscription/internal/logger"
	"online-subscription/internal/model"
	"online-subscription/internal/usecase"
	"strings"
	"time"

	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	uc *usecase.SubscriptionUseCase
}

func NewSubscriptionHandler(uc *usecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{uc: uc}
}

// Create godoc
// @Summary Create a new subscription
// @Description Create a subscription record
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

	helpers.WriteJSON(w, http.StatusCreated, sub)
}

// List godoc
// @Summary List subscriptions
// @Description Get a list of subscriptions with optional filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "Filter by User ID"
// @Param service_name query string false "Filter by Service Name"
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
	helpers.WriteJSON(w, http.StatusOK, subs)
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
	helpers.WriteJSON(w, http.StatusOK, s)
}

// Update godoc
// @Summary Update subscription
// @Description Update fields of an existing subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param body body dto.UpdateSubscriptionRequest true "Fields to update"
// @Success 200 {object} model.Subscription
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [patch]
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := h.uc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sub == nil {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	}

	if req.ServiceName != nil {
		sub.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		sub.Price = *req.Price
	}
	if req.StartDate != nil {
		start, err := helpers.ParseDateToTime(*req.StartDate)
		if err != nil {
			http.Error(w, "invalid start_date format", http.StatusBadRequest)
			return
		}
		sub.StartDate = start
	}
	if req.EndDate != nil && *req.EndDate != "" {
		end, err := helpers.ParseDateToTime(*req.EndDate)
		if err != nil {
			http.Error(w, "invalid end_date format", http.StatusBadRequest)
			return
		}
		sub.EndDate = &end
	} else {
		sub.EndDate = nil
	}

	if err := h.uc.Update(r.Context(), sub); err != nil {
		logger.Error("Failed to update subscription", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Subscription updated",
		zap.String("id", sub.ID),
		zap.String("service", sub.ServiceName),
		zap.String("user_id", sub.UserID),
	)

	helpers.WriteJSON(w, http.StatusOK, sub)
}

// Delete godoc
// @Summary Delete subscription
// @Description Delete a subscription by ID
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
// @Summary Get subscriptions summary
// @Description Calculate total subscription cost for a period with optional filters
// @Tags subscriptions
// @Produce json
// @Param from query string true "Start date in MM-YYYY"
// @Param to query string false "End date in MM-YYYY"
// @Param user_id query string false "Filter by User ID"
// @Param service_name query string false "Filter by Service Name"
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

	fromDate, err := helpers.ParseDateToTime(from)
	if err != nil {
		http.Error(w, "invalid from date", http.StatusBadRequest)
		return
	}

	var toDate *time.Time
	if strings.TrimSpace(to) != "" {
		t, err := helpers.ParseDateToTime(to)
		if err != nil {
			http.Error(w, "invalid to date", http.StatusBadRequest)
			return
		}
		if t.Before(fromDate) {
			http.Error(w, "`to` date cannot be earlier than `from` date", http.StatusBadRequest)
			return
		}
		toDate = &t
	}

	f := model.SummaryFilter{
		FromDate:    fromDate,
		ToDate:      toDate,
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
		zap.String("from", fromDate.Format("01-2006")),
		zap.String("to", func() string {
			if toDate != nil {
				return toDate.Format("01-2006")
			}
			return ""
		}()),
	)

	helpers.WriteJSON(w, http.StatusOK, map[string]int{"total": sum})
}
