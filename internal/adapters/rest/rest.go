package rest

type Handler struct {
	User    *UserHandler
	Manager *ManagerDashboardHandler
	Auth    *AuthHandler
	Product *ProductHandler
	Staff   *StaffHandler
}

func ProvideHandler(
	user *UserHandler,
	manager *ManagerDashboardHandler,
	auth *AuthHandler,
	product *ProductHandler,
	staff *StaffHandler,
) *Handler {
	return &Handler{
		User:    user,
		Manager: manager,
		Auth:    auth,
		Product: product,
		Staff:   staff,
	}
}