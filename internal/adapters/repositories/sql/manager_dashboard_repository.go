package sql

import (
	"context"
	"database/sql"
	"sort"
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
    d, err := r.q.RevenueDurations(ctx, dbmodel.RevenueDurationsParams{Start: start, End: end})
    if err != nil { return 0, err }
    s, err := r.q.RevenueSessions(ctx, dbmodel.RevenueSessionsParams{Start: start, End: end})
    if err != nil { return 0, err }
    return d + s, nil
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
    drows, err := r.q.TopSellingProductsDurations(ctx, dbmodel.TopSellingProductsDurationsParams{Start: start, End: end})
    if err != nil { return nil, err }
    srows, err := r.q.TopSellingProductsSessions(ctx, dbmodel.TopSellingProductsSessionsParams{Start: start, End: end})
    if err != nil { return nil, err }

    m := map[string]int64{}
    for _, r := range drows { m[r.Name] += r.Units }
    for _, r := range srows { m[r.Name] += r.Units }

    out := make([]domrepo.TopProduct, 0, len(m))
    for name, units := range m { out = append(out, domrepo.TopProduct{Name: name, Units: units}) }

    sort.Slice(out, func(i, j int) bool { return out[i].Units > out[j].Units })
    if len(out) > 5 { out = out[:5] }
    return out, nil
}