package requests

import (
	"time"
)

type ManagerDashboardRequest struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}