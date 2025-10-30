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
	Booking          *BookingHandler
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
	booking *BookingHandler,
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
		Booking:          booking,
	}
}
