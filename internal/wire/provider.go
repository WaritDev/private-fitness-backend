package wire

import (
	"github.com/WaritDev/private-fitness-backend/config"
	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/repositories/sql"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/rest"
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

	sql.ProvideTrainerRepository,        // Already returns interface, no need to bind
	sql.ProvidePaymentAccountRepository, // Already returns interface, no need to bind
	// Bind adapters -> domain interfaces
	wire.Bind(new(repositories.UserRepo), new(*sql.UserRepository)),
	wire.Bind(new(repositories.AuthRepo), new(*sql.AuthRepository)),
	wire.Bind(new(repositories.StaffRepository), new(*sql.StaffRepository)),
)

var ServiceSet = wire.NewSet(
	usecases.ProvideUserUseCase,
	usecases.ProvideAuthUseCase,
	usecases.ProvideManagerDashboardUsecase,
	usecases.ProvideProductUseCase,
	usecases.ProvideStaffUsecase,
	usecases.ProvideSessionUseCase,
	usecases.ProvidePaymentUseCase,
)

var HandlerSet = wire.NewSet(
	rest.ProvideUserHandler,
	rest.ProvideAuthHandler,
	rest.ProvideHandler,
	rest.ProvideManagerDashboardHandler,
	rest.ProvideProductHandler,
	rest.ProvideStaffHandler,
	rest.ProvideTrainerHandler,
	rest.ProvidePaymentHandler,
)
