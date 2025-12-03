package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ParseDateToTime(str string) (time.Time, error) {
	str = strings.TrimSpace(str)

	t, err := time.Parse("01-2006", str)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected MM-YYYY")
	}

	return t, nil
}

func PtrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
