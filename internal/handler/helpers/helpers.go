package helpers

import (
	"fmt"
	"strings"
	"time"
)

func ParseDate(str string) (int, int, error) {
	str = strings.TrimSpace(str)
	t, err := time.Parse("01-2006", str)
	if err != nil {

		return 0, 0, fmt.Errorf("invalid date format, expected MM-YYYY")
	}

	return int(t.Month()), t.Year(), nil
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
