package services

import (
	"context"
	"log"
	"strings"

	"rule-based-approval-engine/internal/app/repositories"
	"rule-based-approval-engine/internal/constants"
	"rule-based-approval-engine/internal/models"
	"rule-based-approval-engine/internal/pkg/apperrors"
	"rule-based-approval-engine/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	AdminEmail   = "admin@company.com"
	ManagerEmail = "manager@company.com"

	AdminID   = int64(1)
	ManagerID = int64(2)
)

// AuthService handles authentication and user registration business logic
type AuthService struct {
	userRepo    repositories.UserRepository
	balanceRepo repositories.BalanceRepository
	db          *pgxpool.Pool
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repositories.UserRepository, balanceRepo repositories.BalanceRepository, db *pgxpool.Pool) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		balanceRepo: balanceRepo,
		db:          db,
	}
}

// RegisterUser registers a new user
func (s *AuthService) RegisterUser(ctx context.Context, name, email, password string) error {
	log.Println("RegisterUser started:", email)

	if strings.TrimSpace(email) == "" {
		log.Println("Validation failed: email empty")
		return apperrors.ErrEmailRequired
	}
	if strings.TrimSpace(password) == "" {
		log.Println("Validation failed: password empty")
		return apperrors.ErrPasswordRequired
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Println("DB Begin failed:", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Check email uniqueness
	exists, err := s.userRepo.CheckEmailExists(ctx, tx, email)
	if err != nil {
		log.Println("Email uniqueness query failed:", err)
		return err
	}
	if exists {
		log.Println("Email already registered:", email)
		return apperrors.ErrEmailAlreadyRegistered
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Password hashing failed:", err)
		return apperrors.ErrPasswordHashFailed
	}

	// Decide role, grade, manager_id
	var role string
	var gradeID int64
	var managerID *int64

	switch email {
	case AdminEmail:
		role = constants.RoleAdmin
		gradeID = 3
		managerID = nil

	case ManagerEmail:
		role = constants.RoleManager
		gradeID = 2
		managerID = &AdminID

	default:
		role = constants.RoleEmployee
		gradeID = 1
		managerID = &ManagerID
	}

	log.Printf("Role decided: role=%s grade=%d managerID=%v\n", role, gradeID, managerID)

	// Insert user
	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		GradeID:      gradeID,
		Role:         role,
		ManagerID:    managerID,
	}

	userID, err := s.userRepo.Create(ctx, tx, user)
	if err != nil {
		log.Println("User insert failed:", err)
		return err
	}

	log.Println("User inserted successfully, userID:", userID)

	// Initialize balances ONLY for employee & manager
	if role != constants.RoleAdmin {
		log.Println("Initializing balances for user:", userID)

		err = s.balanceRepo.InitializeBalances(ctx, tx, userID, gradeID)
		if err != nil {
			log.Println("InitializeBalances failed:", err)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		return err
	}

	log.Println("RegisterUser completed successfully:", email)
	return nil
}

// LoginUser authenticates a user and returns a JWT token
func (s *AuthService) LoginUser(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", apperrors.ErrInvalidCredentials
	}

	if err := utils.CheckPassword(password, user.PasswordHash); err != nil {
		return "", "", apperrors.ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	return token, user.Role, nil
}
