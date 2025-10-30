package usecases

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type CustomerLogUsecase struct {
	repo repositories.CustomerLogRepository
}

func ProvideCustomerLogUsecase(repo repositories.CustomerLogRepository) *CustomerLogUsecase {
	return &CustomerLogUsecase{repo: repo}
}

func (uc *CustomerLogUsecase) List(ctx context.Context, req requests.ListCustomerLogsRequest) (responses.ListCustomerLogsResponse, error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	data, err := uc.repo.List(ctx, limit, offset)
	if err != nil {
		return responses.ListCustomerLogsResponse{}, err
	}

	if len(data) == 0 {
		data = []dbmodel.ListCustomerLogsRow{}
	}

	total, err := uc.repo.Count(ctx)
	if err != nil {
		return responses.ListCustomerLogsResponse{}, err
	}

	return responses.ListCustomerLogsResponse{
		Data: data,
		Meta: responses.PageMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: int32(math.Ceil(float64(total) / float64(limit))),
		},
	}, nil
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