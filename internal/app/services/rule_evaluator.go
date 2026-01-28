package services

import "fmt"

func EvaluateLeaveRule(rule map[string]interface{}, days int) bool {
	maxDays, ok := rule["max_days"].(float64)
	fmt.Println("EvaluateLeaveRule ", ok)
	if !ok {
		return false
	}
	return days <= int(maxDays)
}

func EvaluateExpenseRule(rule map[string]interface{}, amount float64) bool {
	maxAmount, ok := rule["max_amount"].(float64)
	if !ok {
		return false
	}
	return amount <= maxAmount
}

func EvaluateDiscountRule(rule map[string]interface{}, percent float64) bool {
	maxPercent, ok := rule["max_percent"].(float64)
	if !ok {
		return false
	}
	return percent <= maxPercent
}
