package rest
type Handler struct {
	User    *UserHandler
	Manager *ManagerDashboardHandler
	Auth *AuthHandler
	Product *ProductHandler
}

func ProvideHandler(
	user *UserHandler,
	manager *ManagerDashboardHandler,
	auth *AuthHandler,
	product *ProductHandler,
) *Handler {
	return &Handler{
		User:            user,
		Manager:         manager,
		Auth:           auth,
		Product:        product,
	}
}