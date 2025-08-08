//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"BlockAccount/internal/client"
	"BlockAccount/internal/conf"
	"BlockAccount/internal/data"
	"BlockAccount/internal/handler"
	"BlockAccount/internal/server"
	"BlockAccount/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, handler.ProviderSet, service.ProviderSet, client.ProviderSet, newApp))
}
