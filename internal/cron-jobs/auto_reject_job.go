package jobs

import (
	"log"
	"rule-based-approval-engine/internal/app/services"
)

func RunAutoRejectJob() {
	log.Println("⏱️ Auto-reject job started")

	services.AutoRejectLeaveRequests()
	services.AutoRejectExpenseRequests()
	services.AutoRejectDiscountRequests()

	log.Println("✅ Auto-reject job finished")
}
