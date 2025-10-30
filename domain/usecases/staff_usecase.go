package usecases

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type StaffUsecase struct {
	repo repositories.StaffRepository
}

func ProvideStaffUsecase(repo repositories.StaffRepository) *StaffUsecase {
	return &StaffUsecase{repo: repo}
}

func (uc *StaffUsecase) ListStaffs(ctx context.Context, req requests.ListStaffsRequest) (responses.ListStaffsResponse, error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 { limit = 10 }
	page := req.Page
	if page <= 0 { page = 1 }
	offset := (page - 1) * limit

	data, err := uc.repo.List(ctx, limit, offset)
	if err != nil { return responses.ListStaffsResponse{}, err }
	total, err := uc.repo.Count(ctx)
	if err != nil { return responses.ListStaffsResponse{}, err }

	return responses.ListStaffsResponse{
		Data: data,
		Meta: responses.PageMeta{
			Page: page, Limit: limit,
			TotalItems: total,
			TotalPages: int32(math.Ceil(float64(total)/float64(limit))),
		},
	}, nil
}

func (uc *StaffUsecase) CreateStaff(ctx context.Context, req requests.CreateStaffRequest) (responses.StaffCreatedResponse, error) {
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

func (uc *StaffUsecase) UpdateStaff(ctx context.Context, targetUsername string, req requests.UpdateStaffRequest) (responses.StaffUpdatedResponse, error) {
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

func (uc *StaffUsecase) DeleteStaff(ctx context.Context, targetUsername string) (responses.StaffDeletedResponse, error) {
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