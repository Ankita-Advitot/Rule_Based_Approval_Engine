package services

import (
	"context"
	"errors"
	"log"
	"rule-based-approval-engine/internal/database"
	"rule-based-approval-engine/internal/models"
	"rule-based-approval-engine/internal/pkg/utils"
	"strings"
)

var (
	AdminEmail   = "admin@company.com"
	ManagerEmail = "manager@company.com"

	AdminID   = int64(1)
	ManagerID = int64(2)
)

func RegisterUser(name, email, password string) error {
	ctx := context.Background()

	log.Println("RegisterUser started:", email)
	if strings.TrimSpace(email) == "" {
		log.Println("Validation failed: email empty")
		return errors.New("email is required")
	}
	if strings.TrimSpace(password) == "" {
		log.Println("Validation failed: password empty")
		return errors.New("password is required")
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		log.Println("DB Begin failed:", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Check email uniqueness
	var count int
	err = tx.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM users WHERE email=$1",
		email,
	).Scan(&count)

	if err != nil {
		log.Println("Email uniqueness query failed:", err)
		return err
	}
	if count > 0 {
		log.Println("Email already registered:", email)
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Password hashing failed:", err)
		return err
	}

	// Decide role, grade, manager_id
	var role string
	var gradeID int64
	var managerID *int64

	switch email {
	case AdminEmail:
		role = "ADMIN"
		gradeID = 3
		managerID = nil

	case ManagerEmail:
		role = "MANAGER"
		gradeID = 2
		managerID = &AdminID

	default:
		role = "EMPLOYEE"
		gradeID = 1
		managerID = &ManagerID
	}

	log.Printf("Role decided: role=%s grade=%d managerID=%v\n", role, gradeID, managerID)

	// Insert user
	var userID int64
	err = tx.QueryRow(
		ctx,
		`INSERT INTO users (name, email, password_hash, grade_id, role, manager_id)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id`,
		name, email, hashedPassword, gradeID, role, managerID,
	).Scan(&userID)

	if err != nil {
		log.Println("User insert failed:", err)
		return err
	}

	log.Println("User inserted successfully, userID:", userID)

	// Initialize balances ONLY for employee & manager
	if role != "ADMIN" {
		log.Println("Initializing balances for user:", userID)

		err = InitializeBalances(tx, userID, gradeID)
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

func LoginUser(email, password string) (string, string, error) {
	var user models.User

	err := database.DB.QueryRow(
		context.Background(),
		`SELECT id, password_hash, role FROM users WHERE email=$1`,
		email,
	).Scan(&user.ID, &user.PasswordHash, &user.Role)

	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := utils.CheckPassword(password, user.PasswordHash); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	return token, user.Role, nil
}
