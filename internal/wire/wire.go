//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/rest"
)

func InitializeHandler() *rest.Handler {
	wire.Build(
		ConfigSet,
		InfraSet,
		ProvideContext,
		RepositorySet,
		ServiceSet,
		HandlerSet,
	)
	return &rest.Handler{}
}