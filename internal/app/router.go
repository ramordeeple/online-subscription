package app

import (
	"net/http"
	"online-subscription/internal/handler"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handler.SubscriptionHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/subscriptions/summary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.Summary(w, r)
	})

	mux.HandleFunc("/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.List(w, r)
		case http.MethodPost:
			h.Create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/subscriptions/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.GetById(w, r, id)
		case http.MethodDelete:
			h.Delete(w, r, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	return mux
}
