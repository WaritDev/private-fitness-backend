package usecases

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

type CustomerUsecase struct {
	repo repositories.CustomerRepository
}

func ProvideCustomerUsecase(repo repositories.CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{repo: repo}
}

func (uc *CustomerUsecase) ListCustomers(ctx context.Context) ([]dbmodel.ListCustomersRow, error) {
	return uc.repo.List(ctx)
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

func (uc *CustomerUsecase) GetCustomerByUsername(ctx context.Context, username string) (responses.Customer, error) {
	if strings.TrimSpace(username) == "" {
		return responses.Customer{}, errors.New("username required")
	}

	row, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return responses.Customer{}, err
	}

	var dob string
	switch v := any(row.DateOfBirth).(type) {
	case time.Time:
		dob = v.Format("2006-01-02")
	case *time.Time:
		if v != nil {
			dob = v.Format("2006-01-02")
		}
	case sql.NullTime:
		if v.Valid {
			dob = v.Time.Format("2006-01-02")
		}
	case *sql.NullTime:
		if v != nil && v.Valid {
			dob = v.Time.Format("2006-01-02")
		}
	}

	resp := responses.Customer{
		Username:    row.Username,
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		Gender:      strings.ToUpper(string(row.Gender)),
		DateOfBirth: dob,
		PhoneNumber: row.PhoneNumber,
		Gmail:       row.Gmail,
		IsActive: responses.NullBool{
			Bool:  row.IsActive.Bool,
			Valid: row.IsActive.Valid,
		},
		HealthInfo:                  ns(row.HealthInfo),
		Address:                     ns(row.Address),
		CompanyName:                 ns(row.CompanyName),
		CompanyPosition:             ns(row.CompanyPosition),
		MaritalStatus:               ns(string(row.MaritalStatus)),
		EmergencyContactName:        ns(row.EmergencyContactName),
		EmergencyContactRelationship: ns(row.EmergencyContactRelationship),
		EmergencyContactPhone:       ns(row.EmergencyContactPhone),
		MarketingSource:             ns(row.MarketingSource),
	}

	return resp, nil
}

func ns(s string) responses.NullString {
	return responses.NullString{String: s, Valid: s != ""}
}