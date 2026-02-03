package services

import (
	"context"
	"time"

	"rule-based-approval-engine/internal/app/repositories"
	"rule-based-approval-engine/internal/constants"
	"rule-based-approval-engine/internal/pkg/apperrors"
)

type HolidayService struct {
	holidayRepo repositories.HolidayRepository
}

func NewHolidayService(holidayRepo repositories.HolidayRepository) *HolidayService {
	return &HolidayService{holidayRepo: holidayRepo}
}

func (s *HolidayService) ensureAdmin(role string) error {
	if role != constants.RoleAdmin {
		return apperrors.ErrAdminOnly
	}
	return nil
}

func (s *HolidayService) AddHoliday(ctx context.Context, role string, adminID int64, date time.Time, desc string) error {
	if err := s.ensureAdmin(role); err != nil {
		return err
	}
	return s.holidayRepo.AddHoliday(ctx, date, desc, adminID)
}

func (s *HolidayService) GetHolidays(ctx context.Context, role string) ([]map[string]interface{}, error) {
	if err := s.ensureAdmin(role); err != nil {
		return nil, err
	}
	return s.holidayRepo.GetHolidays(ctx)
}

func (s *HolidayService) DeleteHoliday(ctx context.Context, role string, holidayID int64) error {
	if err := s.ensureAdmin(role); err != nil {
		return err
	}
	return s.holidayRepo.DeleteHoliday(ctx, holidayID)
}
