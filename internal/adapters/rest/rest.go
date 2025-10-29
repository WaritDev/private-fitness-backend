package rest
type Handler struct {
	User    *UserHandler
	Manager *ManagerDashboardHandler
}

func ProvideHandler(
	user *UserHandler,
	manager *ManagerDashboardHandler,
) *Handler {
	return &Handler{
		User:    user,
		Manager: manager,
	}
}