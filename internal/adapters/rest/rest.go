package rest

type Handler struct {
	User             *UserHandler
	Manager          *ManagerDashboardHandler
	Auth             *AuthHandler
	Product          *ProductHandler
	Trainer          *TrainerHandler
	Payment          *PaymentHandler
	CustomerSession  *CustomerSessionHandler
	CustomerDuration *CustomerDurationHandler
}

func ProvideHandler(
	user *UserHandler,
	manager *ManagerDashboardHandler,
	auth *AuthHandler,
	product *ProductHandler,
	trainer *TrainerHandler,
	payment *PaymentHandler,
	customerSession *CustomerSessionHandler,
	customerDuration *CustomerDurationHandler,
) *Handler {
	return &Handler{
		User:             user,
		Manager:          manager,
		Auth:             auth,
		Product:          product,
		Trainer:          trainer,
		Payment:          payment,
		CustomerSession:  customerSession,
		CustomerDuration: customerDuration,
	}
}
