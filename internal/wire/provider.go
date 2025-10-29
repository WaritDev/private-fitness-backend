package wire

import (
	"github.com/WaritDev/private-fitness-backend/config"
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

var ServiceSet = wire.NewSet(
	usecases.ProvideUserUseCase,
	usecases.ProvideManagerDashboardUsecase,
)

var RepositorySet = wire.NewSet(
	sql.ProvideUserRepository,
	sql.ProvideManagerDashboardRepository,
)

var HandlerSet = wire.NewSet(
	rest.ProvideUserHandler,
	rest.ProvideHandler,
	rest.ProvideManagerDashboardHandler,
)