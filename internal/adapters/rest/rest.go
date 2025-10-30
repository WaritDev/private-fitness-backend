package rest

type Handler struct {
	User             *UserHandler
	Manager          *ManagerDashboardHandler
	Auth             *AuthHandler
	Product          *ProductHandler
	Staff            *StaffHandler
	Trainer          *TrainerHandler
	Payment          *PaymentHandler
	Customer         *CustomerHandler
	CustomerSession  *CustomerSessionHandler
	CustomerDuration *CustomerDurationHandler
	Booking          *BookingHandler
	Member           *MemberHandler
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
	customerSession *CustomerSessionHandler,
	customerDuration *CustomerDurationHandler,
	booking *BookingHandler,
	member *MemberHandler,
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
		Staff:            staff,
		Customer:         customer,
		Member:           member,
	}
}
