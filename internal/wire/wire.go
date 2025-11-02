//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/rest/controller"
)

func InitializeController() *controller.Controller {
	wire.Build(
		ConfigSet,
		InfraSet,
		ProvideContext,
		RepositorySet,
		ServiceSet,
		ControllerSet,
	)
	return &controller.Controller{}
}