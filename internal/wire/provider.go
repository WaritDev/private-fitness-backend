package wire

import (
	"github.com/WaritDev/private-fitness-backend/config"
	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/services"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/repositories/sql"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/rest/controller"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db"

	"github.com/google/wire"
)

var ConfigSet = wire.NewSet(
	config.ProvideConfig,
)

var InfraSet = wire.NewSet(
	db.ProvideMariaDB,
	db.ProvideQueries,
)

var RepositorySet = wire.NewSet(
	sql.ProvideUserRepository,
	sql.ProvideAuthRepository,
	sql.ProvideManagerDashboardRepository,
	sql.ProvideProductRepository,
	sql.ProvideStaffRepository,
	sql.ProvideCustomerRepository,

	sql.ProvideTrainerRepository,          // Already returns interface, no need to bind
	sql.ProvidePaymentAccountRepository,   // Already returns interface, no need to bind
	sql.ProvideCustomerSessionRepository,  // Already returns interface, no need to bind
	sql.ProvideCustomerDurationRepository, // Already returns interface, no need to bind
	sql.ProvideTrainingScheduleRepository, // Already returns interface, no need to bind
	sql.ProvideCustomerLogRepository,      // Already returns interface, no need to bind
	// Bind adapters -> domain interfaces
	wire.Bind(new(repositories.UserRepo), new(*sql.UserRepository)),
	wire.Bind(new(repositories.AuthRepo), new(*sql.AuthRepository)),
	wire.Bind(new(repositories.StaffRepository), new(*sql.StaffRepository)),
	wire.Bind(new(repositories.CustomerRepository), new(*sql.CustomerRepository)),
)

var ServiceSet = wire.NewSet(
	services.ProvideUserService,
	services.ProvideAuthService,
	services.ProvideManagerDashboardService,
	services.ProvideProductService,
	services.ProvideStaffService,
	services.ProvideSessionService,
	services.ProvidePaymentService,
	services.ProvideCustomerService,
	services.ProvideCustomerDurationService,
	services.ProvideCustomerSessionService,
	services.ProvideBookingService,
	services.ProvideCustomerLogService,
	services.ProvideMemberService,
	services.ProvideTrainerService, // Use Case 1P: Manage Working Hours
)

var ControllerSet = wire.NewSet(
	controller.ProvideUserController,
	controller.ProvideAuthController,
	controller.ProvideController,
	controller.ProvideManagerDashboardController,
	controller.ProvideProductController,
	controller.ProvideStaffController,
	controller.ProvideTrainerController,
	controller.ProvidePaymentController,
	controller.ProvideCustomerDurationController,
	controller.ProvideCustomerController,
	controller.ProvideCustomerSessionController,
	controller.ProvideBookingController,
	controller.ProvideCustomerLogController,
	controller.ProvideMemberController,
)
