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

type CustomerUsecase struct {
	repo repositories.CustomerRepository
}

func ProvideCustomerUsecase(repo repositories.CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{repo: repo}
}

func (uc *CustomerUsecase) ListCustomers(ctx context.Context, req requests.ListCustomersRequest) (responses.ListCustomersResponse, error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 { limit = 10 }
	page := req.Page
	if page <= 0 { page = 1 }
	offset := (page - 1) * limit

	data, err := uc.repo.List(ctx, limit, offset)
	if err != nil { return responses.ListCustomersResponse{}, err }

	total, err := uc.repo.Count(ctx)
	if err != nil { return responses.ListCustomersResponse{}, err }

	return responses.ListCustomersResponse{
		Data: data,
		Meta: responses.PageMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: int32(math.Ceil(float64(total)/float64(limit))),
		},
	}, nil
}

func (uc *CustomerUsecase) UpdateCustomer(ctx context.Context, targetUsername string, req requests.UpdateCustomerRequest) (responses.CustomerUpdatedResponse, error) {
	// validate พื้นฐาน
	if strings.TrimSpace(targetUsername) == "" {
		return responses.CustomerUpdatedResponse{}, errors.New("username required")
	}
	if strings.TrimSpace(req.FirstName) == "" { return responses.CustomerUpdatedResponse{}, errors.New("firstName required") }
	if strings.TrimSpace(req.LastName)  == "" { return responses.CustomerUpdatedResponse{}, errors.New("lastName required") }

	gender := strings.ToUpper(strings.TrimSpace(req.Gender))
	switch gender {
	case "MALE","FEMALE","OTHER":
	default: return responses.CustomerUpdatedResponse{}, errors.New("invalid gender")
	}

	if !utils.IsValidPhone(req.PhoneNumber) {
		return responses.CustomerUpdatedResponse{}, errors.New("invalid phone")
	}
	email := utils.NormalizeEmail(req.Gmail)
	if !utils.IsValidEmail(email) {
		return responses.CustomerUpdatedResponse{}, errors.New("invalid gmail")
	}
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil { return responses.CustomerUpdatedResponse{}, errors.New("invalid date_of_birth format (YYYY-MM-DD)") }
	if age(dob) < 14 { return responses.CustomerUpdatedResponse{}, errors.New("age must be >= 14") }

	// ตรวจซ้ำ ยกเว้นเจ้าของ
	if n, err := uc.repo.CheckPhoneExistsExcept(ctx, req.PhoneNumber, targetUsername); err != nil {
		return responses.CustomerUpdatedResponse{}, err
	} else if n > 0 {
		return responses.CustomerUpdatedResponse{}, errors.New("phone already exists")
	}
	if n, err := uc.repo.CheckGmailExistsExcept(ctx, email, targetUsername); err != nil {
		return responses.CustomerUpdatedResponse{}, err
	} else if n > 0 {
		return responses.CustomerUpdatedResponse{}, errors.New("gmail already exists")
	}

	// อัปเดต users (มี/ไม่มี password)
	newPass := strings.TrimSpace(getString(req.NewPassword))
	if newPass == "" {
		if err := uc.repo.UpdateUserNoPassword(ctx, repositories.UpdateCustomerUserParams{
			Username:    targetUsername,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Gender:      gender,
			DateOfBirth: dob,
			PhoneNumber: req.PhoneNumber,
			Gmail:       email,
			IsActive:    req.IsActive,
		}); err != nil {
			return responses.CustomerUpdatedResponse{}, err
		}
	} else {
		if req.ConfirmNewPassword == nil || *req.ConfirmNewPassword != newPass {
			return responses.CustomerUpdatedResponse{}, errors.New("passwords do not match")
		}
		if !utils.ValidatePassword(newPass) {
			return responses.CustomerUpdatedResponse{}, errors.New("weak password")
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
		if err != nil { return responses.CustomerUpdatedResponse{}, err }

		if err := uc.repo.UpdateUserWithPassword(ctx, repositories.UpdateCustomerUserWithPassParams{
			UpdateCustomerUserParams: repositories.UpdateCustomerUserParams{
				Username:    targetUsername,
				FirstName:   req.FirstName,
				LastName:    req.LastName,
				Gender:      gender,
				DateOfBirth: dob,
				PhoneNumber: req.PhoneNumber,
				Gmail:       email,
				IsActive:    req.IsActive,
			},
			HashedPass: string(hashed),
		}); err != nil {
			return responses.CustomerUpdatedResponse{}, err
		}
	}

	// อัปเดตตาราง customers
	if err := uc.repo.UpdateCustomerDetail(ctx, repositories.UpdateCustomerDetailParams{
		Username:                     targetUsername,
		HealthInfo:                   req.HealthInfo,
		Address:                      req.Address,
		CompanyName:                  req.CompanyName,
		CompanyPosition:              req.CompanyPosition,
		MaritalStatus:                req.MaritalStatus,
		EmergencyContactName:         req.EmergencyContactName,
		EmergencyContactRelationship: req.EmergencyContactRelationship,
		EmergencyContactPhone:        req.EmergencyContactPhone,
		MarketingSource:              req.MarketingSource,
	}); err != nil {
		return responses.CustomerUpdatedResponse{}, err
	}

	return responses.CustomerUpdatedResponse{Message: "Customer: " + targetUsername + " updated successfully"}, nil
}

func (uc *CustomerUsecase) DeleteCustomer(ctx context.Context, targetUsername string) (responses.CustomerDeletedResponse, error) {
	if strings.TrimSpace(targetUsername) == "" {
		return responses.CustomerDeletedResponse{}, errors.New("username required")
	}
	if err := uc.repo.DeleteCustomerCascade(ctx, targetUsername); err != nil {
		return responses.CustomerDeletedResponse{}, err
	}
	return responses.CustomerDeletedResponse{
		Message: "Customer: " + targetUsername + " deleted successfully",
	}, nil
}