package controller

type Controller struct {
	User             *UserController
	Manager          *ManagerDashboardController
	Auth             *AuthController
	Product          *ProductController
	Staff            *StaffController
	Trainer          *TrainerController
	Payment          *PaymentController
	Customer         *CustomerController
	CustomerSession  *CustomerSessionController
	CustomerDuration *CustomerDurationController
	Booking          *BookingController
	Member           *MemberController
	CustomerLog *CustomerLogController
}

func ProvideController(
	user *UserController,
	manager *ManagerDashboardController,
	auth *AuthController,
	product *ProductController,
	trainer *TrainerController,
	payment *PaymentController,
	staff *StaffController,
	customer *CustomerController,
	customerSession *CustomerSessionController,
	customerDuration *CustomerDurationController,
	booking *BookingController,
	member *MemberController,
	customerLog *CustomerLogController,
) *Controller {
	return &Controller{
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
		CustomerLog:      customerLog,
	}
}
