🧩 Go Clean Architecture Refactor Guide

From Usecase → Service and Handler/Rest → Controller

⸻

🎯 Objective

This document defines the standardized naming convention and step-by-step process to refactor the existing Go Clean Architecture project.
The goal is to replace all legacy Usecase and Handler/Rest naming patterns with Service and Controller respectively, improving clarity and aligning with modern domain naming standards.

⸻

🗂️ 1. Overview of Naming Rules

Layer	Old Convention	New Convention	Example
Application Layer	*Usecase, _use_case.go or _usecase.go	*Service, _service.go	CustomerUsecase → CustomerService
Delivery Layer	*Handler, _rest.go	*Controller, _controller.go	CustomerHandler → CustomerController


⸻

🧱 2. Directory Structure (After Refactor)

internal/
├── app/
│   └── service/
│       ├── auth_service.go
│       ├── booking_service.go
│       ├── customer_log_service.go
│       ├── customer_session_service.go
│       ├── customer_service.go
│       ├── duration_service.go
│       ├── manager_dashboard_service.go
│       ├── member_service.go
│       ├── payment_service.go
│       ├── product_service.go
│       ├── session_service.go
│       ├── staff_service.go
│       ├── trainer_service.go
│       └── user_service.go
│
└── adapter/
    └── rest/
        └── controller/
            ├── auth_controller.go
            ├── booking_controller.go
            ├── customer_log_controller.go
            ├── customer_session_controller.go
            ├── customer_controller.go
            ├── duration_controller.go
            ├── manager_dashboard_controller.go
            ├── member_controller.go
            ├── payment_controller.go
            ├── product_controller.go
            ├── session_controller.go
            ├── staff_controller.go
            ├── trainer_controller.go
            └── user_controller.go


⸻

🔁 3. File & Class Rename Mapping

🧩 Application Layer (Usecase → Service)

Module	Old File	New File	Old Struct	New Struct
Auth	auth_use_case.go	auth_service.go	AuthUsecase	AuthService
Booking	booking_use_case.go	booking_service.go	BookingUsecase	BookingService
Customer Log	customer_log_usecase.go	customer_log_service.go	CustomerLogUsecase	CustomerLogService
Customer Session	customer_session_use_case.go	customer_session_service.go	CustomerSessionUsecase	CustomerSessionService
Customer	customer_usecase.go	customer_service.go	CustomerUsecase	CustomerService
Duration	duration_use_case.go	duration_service.go	DurationUsecase	DurationService
Manager Dashboard	manager_dashboard_usecase.go	manager_dashboard_service.go	ManagerDashboardUsecase	ManagerDashboardService
Member	member_use_case.go	member_service.go	MemberUsecase	MemberService
Payment	payment_use_case.go	payment_service.go	PaymentUsecase	PaymentService
Product	product_use_case.go	product_service.go	ProductUsecase	ProductService
Session	session_use_case.go	session_service.go	SessionUsecase	SessionService
Staff	staff_usecase.go	staff_service.go	StaffUsecase	StaffService
Trainer	trainer_use_case.go	trainer_service.go	TrainerUsecase	TrainerService
User	user_use_case.go	user_service.go	UserUsecase	UserService


⸻

🌐 Delivery Layer (Handler / Rest → Controller)

Module	Old File	New File	Old Struct	New Struct
Auth	auth_rest.go	auth_controller.go	AuthHandler	AuthController
Booking	booking_rest.go	booking_controller.go	BookingHandler	BookingController
Customer	customer_controller.go	(already correct)	✅	✅
Customer Log	customer_log_controller.go	(already correct)	✅	✅
Customer Session	customer_session_rest.go	customer_session_controller.go	CustomerSessionHandler	CustomerSessionController
Duration	duration_rest.go	duration_controller.go	DurationHandler	DurationController
Manager Dashboard	manager_dashboard_controller.go	(already correct)	✅	✅
Member	member_rest.go	member_controller.go	MemberHandler	MemberController
Payment	payment_rest.go	payment_controller.go	PaymentHandler	PaymentController
Product	product_rest.go	product_controller.go	ProductHandler	ProductController
Session	session_rest.go	session_controller.go	SessionHandler	SessionController
Staff	staff_controller.go	(already correct)	✅	✅
Trainer	trainer_rest.go	trainer_controller.go	TrainerHandler	TrainerController
User	user_rest.go	user_controller.go	UserHandler	UserController


⸻

🧩 4. Example: Before & After Refactor

Before

// rest/customer_rest.go
type CustomerHandler struct { // old
	uc *usecases.CustomerUsecase
}

func ProvideCustomerHandler(uc *usecases.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{uc: uc}
}

// usecases/customer_usecase.go
type CustomerUsecase struct {
	repo repositories.CustomerRepository
}

func ProvideCustomerUsecase(repo repositories.CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{repo: repo}
}


⸻

After

// controller/customer_controller.go
package controller

import "github.com/yourorg/yourproject/internal/app/service"

type CustomerController struct {
	svc *service.CustomerService
}

func ProvideCustomerController(svc *service.CustomerService) *CustomerController {
	return &CustomerController{svc: svc}
}

// service/customer_service.go
package service

import "github.com/yourorg/yourproject/internal/domain/repositories"

type CustomerService struct {
	repo repositories.CustomerRepository
}

func ProvideCustomerService(repo repositories.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}


⸻

⚙️ 5. Update Dependency Injection (Wire)

Before

var CustomerSet = wire.NewSet(
	repositories.ProvideCustomerRepository,
	usecases.ProvideCustomerUsecase,
	rest.ProvideCustomerHandler,
)

After

var CustomerSet = wire.NewSet(
	repositories.ProvideCustomerRepository,
	service.ProvideCustomerService,
	controller.ProvideCustomerController,
)