package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"github.com/WaritDev/private-fitness-backend/utils"
)

type CustomerLogUsecase struct {
	repo repositories.CustomerLogRepository
}

func ProvideCustomerLogUsecase(repo repositories.CustomerLogRepository) *CustomerLogUsecase {
	return &CustomerLogUsecase{repo: repo}
}

func (uc *CustomerLogUsecase) List(ctx context.Context) ([]dbmodel.ListCustomerLogsRow, error) {
	return uc.repo.List(ctx)
}

func (uc *CustomerLogUsecase) Update(
	ctx context.Context,
	id int32,
	req requests.UpdateCustomerLogRequest,
) (responses.CustomerLogUpdatedResponse, error) {

	// 1) validate timestamp
	ts, err := time.Parse("2006-01-02 15:04:05", req.Timestamp)
	if err != nil {
		return responses.CustomerLogUpdatedResponse{}, errors.New("invalid timestamp format (YYYY-MM-DD HH:MM:SS)")
	}

	// 2) validate enum
	switch req.LogType {
	case "CHECK_IN", "CHECK_OUT", "BOOK_SESSION", "CANCEL_SESSION":
	default:
		return responses.CustomerLogUpdatedResponse{}, errors.New("invalid logType")
	}

	// 3) persist
	affected, err := uc.repo.UpdateByID(ctx, id, ts, req.LogType)
	if err != nil {
		return responses.CustomerLogUpdatedResponse{}, err
	}
	if affected == 0 {
		return responses.CustomerLogUpdatedResponse{}, errors.New("log not found")
	}

	return responses.CustomerLogUpdatedResponse{
		Message: fmt.Sprintf("Log: %d updated successfully", id),
	}, nil
}

func (uc *CustomerLogUsecase) Delete(
	ctx context.Context,
	id int32,
) (responses.CustomerLogDeletedResponse, error) {

	rows, err := uc.repo.DeleteByID(ctx, id)
	if err != nil {
		return responses.CustomerLogDeletedResponse{}, err
	}
	if rows == 0 {
		return responses.CustomerLogDeletedResponse{}, errors.New("log not found")
	}

	return responses.CustomerLogDeletedResponse{
		Message: fmt.Sprintf("Log: %d deleted successfully", id),
	}, nil
}

func (uc *CustomerLogUsecase) GetByID(ctx context.Context, id string) (responses.CustomerLog, error) {
	if strings.TrimSpace(id) == "" {
		return responses.CustomerLog{}, errors.New("id required")
	}

	id32, err := utils.Atoi32(id)
	if err != nil {
		return responses.CustomerLog{}, errors.New("invalid id")
	}

	row, err := uc.repo.GetByID(ctx, id32)
	if err != nil {
		return responses.CustomerLog{}, err
	}

	resp := responses.CustomerLog{
		ID:                utils.Itoa(row.ID),                    // int32 -> string
		CustomerUsername:  utils.NS(row.CustomerUsername),        // sql.NullString -> string
		CustomerFirstName: row.CustomerFirstName,
		CustomerLastName:  row.CustomerLastName, 
		CreatedAt:         utils.NT(row.CreatedAt),               // sql.NullTime -> RFC3339 string
		LogType:           string(row.LogType),                   // enum alias -> string
	}
	return resp, nil
}