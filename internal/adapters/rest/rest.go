package rest
type Handler struct {
	User    *UserHandler
	Manager *ManagerDashboardHandler
	Auth *AuthHandler
}

func ProvideHandler(
	user *UserHandler,
	manager *ManagerDashboardHandler,
	auth *AuthHandler,
) *Handler {
	return &Handler{
		User:    user,
		Manager: manager,Auth: auth,
	}
}