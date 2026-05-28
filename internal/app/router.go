package app

import (
	"net/http"
	"online-subscription/internal/handler"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handler.SubscriptionHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /subscriptions/summary", h.Summary)
	mux.HandleFunc("GET /subscriptions", h.List)
	mux.HandleFunc("POST /subscriptions", h.Create)
	mux.HandleFunc("GET /subscriptions/{id}", h.GetByID)
	mux.HandleFunc("PATCH /subscriptions/{id}", h.Update)
	mux.HandleFunc("PUT /subscriptions/{id}", h.Update)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.Delete)

	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)
	return mux
}
