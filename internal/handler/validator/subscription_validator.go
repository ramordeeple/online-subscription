package validator

import "fmt"

func ValidateSubscriptionDates(startMonth, startYear int, endMonth, endYear *int) error {
	if endMonth != nil && endYear != nil {
		start := startYear*12 + startMonth
		end := *endYear*12 + *endMonth
		if end < start {
			return fmt.Errorf("end_date cannot be earlier than start_date")
		}
	}
	return nil
}
