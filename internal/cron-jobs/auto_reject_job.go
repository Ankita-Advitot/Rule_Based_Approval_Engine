package jobs

import (
	"context"
	"log"
	"rule-based-approval-engine/internal/app/services"
)

func RunAutoRejectJob(svc *services.AutoRejectService) {
	log.Println("Auto-reject job started")

	ctx := context.Background()
	svc.AutoRejectLeaveRequests(ctx)
	svc.AutoRejectExpenseRequests(ctx)
	svc.AutoRejectDiscountRequests(ctx)

	log.Println("Auto-reject job finished")
}
