package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"github.com/WaritDev/private-fitness-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type StaffService struct {
	repo repositories.StaffRepository
}

func ProvideStaffService(repo repositories.StaffRepository) *StaffService {
	return &StaffService{repo: repo}
}

func (uc *StaffService) ListStaffs(ctx context.Context) ([]dbmodel.ListStaffsRow, error) {
	return uc.repo.List(ctx)
}

func (uc *StaffService) CreateStaff(ctx context.Context, req requests.CreateStaffRequest) (responses.StaffCreatedResponse, error) {
	if !utils.IsValidUsername(req.Username) {
		return responses.StaffCreatedResponse{}, errors.New("invalid username")
	}
	if !utils.ValidatePassword(req.Password) {
		return responses.StaffCreatedResponse{}, errors.New("weak password")
	}
	if req.Password != req.ConfirmPass {
		return responses.StaffCreatedResponse{}, errors.New("passwords do not match")
	}

	role := strings.ToUpper(strings.TrimSpace(req.Role))
	switch role {
	case "TRAINER","SALES","MANAGER","ADMIN":
	default:
		return responses.StaffCreatedResponse{}, errors.New("invalid role")
	}

	gender := strings.ToUpper(strings.TrimSpace(req.Gender))
	switch gender {
	case "MALE","FEMALE","OTHER":
	default:
		return responses.StaffCreatedResponse{}, errors.New("invalid gender")
	}

	if strings.TrimSpace(req.FirstName) == "" {
		return responses.StaffCreatedResponse{}, errors.New("firstName required")
	}
	if strings.TrimSpace(req.LastName) == "" {
		return responses.StaffCreatedResponse{}, errors.New("lastName required")
	}
	if !utils.IsValidPhone(req.PhoneNumber) {
		return responses.StaffCreatedResponse{}, errors.New("invalid phone")
	}

	gmail := utils.NormalizeEmail(req.Gmail)
	if !utils.IsValidEmail(gmail) {
		return responses.StaffCreatedResponse{}, errors.New("invalid email")
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return responses.StaffCreatedResponse{}, errors.New("invalid date_of_birth format (YYYY-MM-DD)")
	}
	if age(dob) < 14 {
		return responses.StaffCreatedResponse{}, errors.New("age must be >= 14")
	}

	spec := req.Specialty
	if role != "TRAINER" {
		spec = nil
	}

	if n, err := uc.repo.CheckUsernameExists(ctx, req.Username); err != nil {
		return responses.StaffCreatedResponse{}, err
	} else if n > 0 {
		return responses.StaffCreatedResponse{}, errors.New("username already exists")
	}

	if n, err := uc.repo.CheckPhoneExists(ctx, req.PhoneNumber); err != nil {
		return responses.StaffCreatedResponse{}, err
	} else if n > 0 {
		return responses.StaffCreatedResponse{}, errors.New("phone already exists")
	}

	if n, err := uc.repo.CheckGmailExists(ctx, gmail); err != nil {
		return responses.StaffCreatedResponse{}, err
	} else if n > 0 {
		return responses.StaffCreatedResponse{}, errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil { return responses.StaffCreatedResponse{}, err }

	err = uc.repo.Create(ctx, repositories.CreateStaffParams{
		Username:    req.Username,
		HashedPass:  string(hashed),
		Role:        role,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      gender,
		DateOfBirth: dob,
		PhoneNumber: req.PhoneNumber,
		Gmail:       gmail,
		Specialty:   spec,
	})
	if err != nil { return responses.StaffCreatedResponse{}, err }

	return responses.StaffCreatedResponse{ Message: "User: " + req.Username + " created successfully" }, nil
}

func age(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() { years-- }
	return years
}

func (uc *StaffService) UpdateStaff(ctx context.Context, targetUsername string, req requests.UpdateStaffRequest) (responses.StaffUpdatedResponse, error) {
	// 1) แปลง/ตรวจค่าพื้นฐาน
	role := strings.ToUpper(strings.TrimSpace(req.Role))
	switch role {
	case "TRAINER","SALES","MANAGER","ADMIN":
	default:
		return responses.StaffUpdatedResponse{}, errors.New("invalid role")
	}

	gender := strings.ToUpper(strings.TrimSpace(req.Gender))
	switch gender {
	case "MALE","FEMALE","OTHER":
	default:
		return responses.StaffUpdatedResponse{}, errors.New("invalid gender")
	}

	if strings.TrimSpace(req.FirstName) == "" {
		return responses.StaffUpdatedResponse{}, errors.New("firstName required")
	}
	if strings.TrimSpace(req.LastName) == "" {
		return responses.StaffUpdatedResponse{}, errors.New("lastName required")
	}
	if !utils.IsValidPhone(req.PhoneNumber) {
		return responses.StaffUpdatedResponse{}, errors.New("invalid phone")
	}

	gmail := utils.NormalizeEmail(req.Gmail)
	if !utils.IsValidEmail(gmail) {
		return responses.StaffUpdatedResponse{}, errors.New("invalid gmail")
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return responses.StaffUpdatedResponse{}, errors.New("invalid date_of_birth format (YYYY-MM-DD)")
	}
	if age(dob) < 14 {
		return responses.StaffUpdatedResponse{}, errors.New("age must be >= 14")
	}

	spec := req.Specialty
	if role != "TRAINER" {
		spec = nil
	}

	if n, err := uc.repo.CheckPhoneExistsExceptUsername(ctx, req.PhoneNumber, targetUsername); err != nil {
		return responses.StaffUpdatedResponse{}, err
	} else if n > 0 {
		return responses.StaffUpdatedResponse{}, errors.New("phone already exists")
	}

	if n, err := uc.repo.CheckGmailExistsExceptUsername(ctx, gmail, targetUsername); err != nil {
		return responses.StaffUpdatedResponse{}, err
	} else if n > 0 {
		return responses.StaffUpdatedResponse{}, errors.New("gmail already exists")
	}

	// 3) ตัดสินใจอัปเดตรหัสผ่านหรือไม่
	newPass := strings.TrimSpace(getString(req.NewPassword))
	if newPass == "" {
		// no password change
		if err := uc.repo.UpdateNoPassword(ctx, repositories.UpdateStaffNoPassParams{
			Username:    targetUsername,
			Role:        role,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Gender:      gender,
			DateOfBirth: dob,
			PhoneNumber: req.PhoneNumber,
			Gmail:       gmail,
			Specialty:   spec,
			IsActive:    req.IsActive,
		}); err != nil {
			return responses.StaffUpdatedResponse{}, err
		}
	} else {
		// validate password + confirm
		if !utils.ValidatePassword(newPass) {
			return responses.StaffUpdatedResponse{}, errors.New("weak password")
		}
		if req.ConfirmNewPassword == nil || *req.ConfirmNewPassword != newPass {
			return responses.StaffUpdatedResponse{}, errors.New("passwords do not match")
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
		if err != nil { return responses.StaffUpdatedResponse{}, err }

		if err := uc.repo.UpdateWithPassword(ctx, repositories.UpdateStaffWithPassParams{
			Username:    targetUsername,
			HashedPass:  string(hashed),
			Role:        role,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Gender:      gender,
			DateOfBirth: dob,
			PhoneNumber: req.PhoneNumber,
			Gmail:       gmail,
			Specialty:   spec,
			IsActive:    req.IsActive,
		}); err != nil {
			return responses.StaffUpdatedResponse{}, err
		}
	}

	return responses.StaffUpdatedResponse{Message: "User: " + targetUsername + " updated successfully"}, nil
}

func getString(p *string) string {
	if p == nil { return "" }
	return *p
}

func (uc *StaffService) DeleteStaff(ctx context.Context, targetUsername string) (responses.StaffDeletedResponse, error) {
	if strings.TrimSpace(targetUsername) == "" {
		return responses.StaffDeletedResponse{}, errors.New("username required")
	}
	if err := uc.repo.DeleteStaffCascade(ctx, targetUsername); err != nil {
		return responses.StaffDeletedResponse{}, err
	}
	return responses.StaffDeletedResponse{
		Message: "Username: " + targetUsername + " deleted successfully",
	}, nil
}

func (uc *StaffService) GetStaffByUsername(ctx context.Context, username string) (responses.Staff, error) {
	if strings.TrimSpace(username) == "" {
		return responses.Staff{}, errors.New("username required")
	}

	row, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return responses.Staff{}, err
	}

	// --- แปลงวันที่ให้เป็น YYYY-MM-DD (รองรับ type ต่าง ๆ ที่ sqlc อาจสร้าง) ---
	var ymd string
	switch v := any(row.DateOfBirth).(type) {
	case time.Time:
		ymd = v.Format("2006-01-02")
	case *time.Time:
		if v != nil {
			ymd = v.Format("2006-01-02")
		}
	case sql.NullTime:
		if v.Valid {
			ymd = v.Time.Format("2006-01-02")
		}
	case *sql.NullTime:
		if v != nil && v.Valid {
			ymd = v.Time.Format("2006-01-02")
		}
	}

	resp := responses.Staff{
		Username:    row.Username,
		Role:        strings.ToUpper(string(row.Role)),
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		Gender:      strings.ToUpper(string(row.Gender)),
		DateOfBirth: ymd,
		PhoneNumber: row.PhoneNumber,
		Gmail:       row.Gmail,
		Specialty: responses.NullString{
			String: row.Specialty.String,
			Valid:  row.Specialty.Valid,
		},
		IsActive: responses.NullBool{
			Bool:  row.IsActive.Bool,
			Valid: row.IsActive.Valid,
		},
	}
	return resp, nil
}