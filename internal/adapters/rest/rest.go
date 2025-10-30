package rest


type Handler struct {
	User    *UserHandler
	Manager *ManagerDashboardHandler
	Auth       *AuthHandler
	Product *ProductHandler
	Staff   *StaffHandler
	Trainer *TrainerHandler
	Payment *PaymentHandler
	Customer *CustomerHandler
	Duration *CustomerDurationHandler
}

func ProvideHandler(
	user *UserHandler,
	manager *ManagerDashboardHandler,
	auth *AuthHandler,
	product *ProductHandler,
	trainer *TrainerHandler,
	payment *PaymentHandler,
	staff *StaffHandler,
	customer *CustomerHandler,
	duration *CustomerDurationHandler,
) *Handler {
	return &Handler{
		User:    user,
		Manager: manager,
		Auth:    auth,
		Product: product,
		Trainer: trainer,
		Payment: payment,
		Staff: staff,
		Customer: customer,
		Duration: duration,
	}
}
