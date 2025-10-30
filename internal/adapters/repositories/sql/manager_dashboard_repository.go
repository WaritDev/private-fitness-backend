package sql

import (
	"context"
	"database/sql"
	"time"

	domrepo "github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type managerDashboardRepo struct {
	q *dbmodel.Queries
}

func ProvideManagerDashboardRepository(db *sql.DB) domrepo.ManagerDashboardRepository {
	return &managerDashboardRepo{q: dbmodel.New(db)}
}

func (r *managerDashboardRepo) TotalRevenue(ctx context.Context, start, end time.Time) (int64, error) {
	v, err := r.q.TotalRevenue(ctx, dbmodel.TotalRevenueParams{
		Start: start,
		End:   end,
	})
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *managerDashboardRepo) NewMembers(ctx context.Context, start, end time.Time) (int64, error) {
	v, err := r.q.NewMembersInRange(ctx, dbmodel.NewMembersInRangeParams{
		Start: start,
		End:   end,
	})
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *managerDashboardRepo) ActiveMembersToday(ctx context.Context) (int64, error) {
	v, err := r.q.ActiveMembersToday(ctx)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *managerDashboardRepo) CheckinsToday(ctx context.Context) (int64, error) {
	v, err := r.q.CheckinsToday(ctx)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *managerDashboardRepo) CompletedPT(ctx context.Context, start, end time.Time) (int64, error) {
	v, err := r.q.CompletedPTInRange(ctx, dbmodel.CompletedPTInRangeParams{
		Start: start,
		End:   end,
	})
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *managerDashboardRepo) TopProducts(ctx context.Context, start, end time.Time) ([]domrepo.TopProduct, error) {
	rows, err := r.q.TopSellingProducts(ctx, dbmodel.TopSellingProductsParams{
		Start: start,
		End:   end,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domrepo.TopProduct, 0, len(rows))
	for _, row := range rows {
		out = append(out, domrepo.TopProduct{
			Name:  row.Name,
			Units: row.Units,
		})
	}
	return out, nil
}