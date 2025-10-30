package repositories

import (
	"context"
	"time"
)

type TopProduct struct {
	Name  string
	Units int64
}

type ManagerDashboardRepository interface {
	TotalRevenue(ctx context.Context, start, end time.Time) (int64, error)
	NewMembers(ctx context.Context, start, end time.Time) (int64, error)
	ActiveMembersToday(ctx context.Context) (int64, error)
	CheckinsToday(ctx context.Context) (int64, error)
	CompletedPT(ctx context.Context, start, end time.Time) (int64, error)
	TopProducts(ctx context.Context, start, end time.Time) ([]TopProduct, error)
}