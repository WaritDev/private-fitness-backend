package usecases

import (
	"context"
	"time"

	domrepo "github.com/WaritDev/private-fitness-backend/domain/repositories"
	res "github.com/WaritDev/private-fitness-backend/domain/responses"
)

type ManagerDashboardUsecase interface {
	Get(ctx context.Context, start, end time.Time) (*res.ManagerDashboardResponse, error)
}

type managerDashboardUsecase struct {
	repo domrepo.ManagerDashboardRepository
}

func ProvideManagerDashboardUsecase(r domrepo.ManagerDashboardRepository) ManagerDashboardUsecase {
	return &managerDashboardUsecase{repo: r}
}

func (uc *managerDashboardUsecase) Get(ctx context.Context, start, end time.Time) (*res.ManagerDashboardResponse, error) {
	total, err := uc.repo.TotalRevenue(ctx, start, end)
	if err != nil { return nil, err }
	newMem, err := uc.repo.NewMembers(ctx, start, end)
	if err != nil { return nil, err }
	active, err := uc.repo.ActiveMembersToday(ctx)
	if err != nil { return nil, err }
	checkins, err := uc.repo.CheckinsToday(ctx)
	if err != nil { return nil, err }
	completedPT, err := uc.repo.CompletedPT(ctx, start, end)
	if err != nil { return nil, err }
	top, err := uc.repo.TopProducts(ctx, start, end)
	if err != nil { return nil, err }

	revSpark := res.Spark{Data: []int64{total / 10, total / 8, total / 6, total / 5, total / 4}}
	newSpark := res.Spark{Data: []int64{newMem / 5, newMem / 4, newMem / 3, newMem / 2, newMem}}
	chkSpark := res.Spark{Data: []int64{checkins / 2, checkins}}
	ptSpark := res.Spark{Data: []int64{completedPT / 3, completedPT / 2, completedPT}}

	return &res.ManagerDashboardResponse{
		TotalRevenueTHB: total,
		NewMembers30d:   newMem,
		ActiveMembers:   active,
		CheckinsToday:   checkins,
		CompletedPT30d:  completedPT,
		RevenueSpark:    revSpark,
		NewMembersSpark: newSpark,
		CheckinsSpark:   chkSpark,
		PTSpark:         ptSpark,
		TopProducts:     func() []res.TopProduct {
			out := make([]res.TopProduct, 0, len(top))
			for _, t := range top {
				out = append(out, res.TopProduct{Name: t.Name, Units: t.Units})
			}
			return out
		}(),
	}, nil
}