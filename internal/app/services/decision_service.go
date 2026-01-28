package services

func Decide(
	requestType string,
	rule map[string]interface{},
	value float64,
) string {

	switch requestType {
	case "LEAVE":
		if EvaluateLeaveRule(rule, int(value)) {
			return "AUTO_APPROVE"
		}
	case "EXPENSE":
		if EvaluateExpenseRule(rule, value) {
			return "AUTO_APPROVE"
		}
	case "DISCOUNT":
		if EvaluateDiscountRule(rule, value) {
			return "AUTO_APPROVE"
		}
	}

	return "MANUAL"
}
