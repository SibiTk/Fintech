//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	
	"cache_service/internal/conf"
	"cache_service/internal/data"
	"cache_service/internal/server"
	"cache_service/internal/service"
"cache_service/internal/handler"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data,*conf.InternalConf, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, service.ProviderSet, handler.ProviderSet, newApp))
}
