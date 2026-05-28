package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"online-subscription/internal/handler/dto"
	"online-subscription/internal/handler/helpers"
	"online-subscription/internal/handler/mapper"
	"online-subscription/internal/handler/parser"
	"online-subscription/internal/model"
	"online-subscription/internal/repository"
	"online-subscription/internal/usecase"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	uc  *usecase.SubscriptionUseCase
	log *zap.Logger
}

func NewSubscriptionHandler(uc *usecase.SubscriptionUseCase, log *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{uc: uc, log: log}
}

func errStatus(err error) int {
	if errors.Is(err, repository.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
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

	h.log.Info("Subscription created",
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
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	f := model.SubscriptionFilter{
		UserID:      helpers.PtrString(q.Get("user_id")),
		ServiceName: helpers.PtrString(q.Get("service_name")),
	}

	if s := q.Get("limit"); s != "" {
		limit, err := strconv.Atoi(s)
		if err != nil || limit < 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		f.Limit = &limit
	}

	if s := q.Get("offset"); s != "" {
		offset, err := strconv.Atoi(s)
		if err != nil || offset < 0 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
		f.Offset = &offset
	}

	subs, err := h.uc.List(r.Context(), &f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Info("Subscriptions listed", zap.Int("count", len(subs)))
	helpers.WriteJSON(w, http.StatusOK, subs)
}

// GetByID godoc
// @Summary Get subscription by ID
// @Description Returns a subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	s, err := h.uc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), errStatus(err))
		return
	}

	h.log.Info("Subscription retrieved", zap.String("id", s.ID))
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
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req dto.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := h.uc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), errStatus(err))
		return
	}

	if req.ServiceName != nil {
		sub.ServiceName = *req.ServiceName
	}
	if req.MonthlyPrice != nil {
		sub.MonthlyPrice = *req.MonthlyPrice
	}
	if req.StartDate != nil {
		start, err := helpers.ParseDateToTime(*req.StartDate)
		if err != nil {
			http.Error(w, "invalid start_date format", http.StatusBadRequest)
			return
		}
		sub.StartDate = start
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			sub.EndDate = nil
		} else {
			end, err := helpers.ParseDateToTime(*req.EndDate)
			if err != nil {
				http.Error(w, "invalid end_date format", http.StatusBadRequest)
				return
			}
			sub.EndDate = &end
		}
	}

	if err := h.uc.Update(r.Context(), sub); err != nil {
		h.log.Error("Failed to update subscription", zap.Error(err))
		http.Error(w, err.Error(), errStatus(err))
		return
	}

	h.log.Info("Subscription updated",
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
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.uc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), errStatus(err))
		return
	}

	h.log.Info("Subscription deleted", zap.String("id", id))
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
	q := r.URL.Query()

	fromDate, err := helpers.ParseDateToTime(q.Get("from"))
	if err != nil {
		http.Error(w, "invalid from date: expected MM-YYYY", http.StatusBadRequest)
		return
	}

	var toDate *time.Time
	if to := strings.TrimSpace(q.Get("to")); to != "" {
		t, err := helpers.ParseDateToTime(to)
		if err != nil {
			http.Error(w, "invalid to date: expected MM-YYYY", http.StatusBadRequest)
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
		UserID:      helpers.PtrString(q.Get("user_id")),
		ServiceName: helpers.PtrString(q.Get("service_name")),
	}

	sum, err := h.uc.Sum(r.Context(), &f)
	if err != nil {
		h.log.Error("Failed to calculate summary", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	toStr := ""
	if toDate != nil {
		toStr = toDate.Format("01-2006")
	}
	h.log.Info("Summary calculated",
		zap.Int("sum", sum),
		zap.String("user_id", helpers.SafeString(f.UserID)),
		zap.String("service_name", helpers.SafeString(f.ServiceName)),
		zap.String("from", fromDate.Format("01-2006")),
		zap.String("to", toStr),
	)

	helpers.WriteJSON(w, http.StatusOK, map[string]int{"total": sum})
}
