//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-mnemosyne-api/user"
)

var userFeatureSet = wire.NewSet(
	user.NewHandler,
	wire.Bind(new(user.Controller), new(*user.Handler)),
	user.NewService,
	wire.Bind(new(user.Service), new(*user.ServiceImpl)),
	user.NewRepository,
	wire.Bind(new(user.Repository), new(*user.RepositoryImpl)),
)

func InitializeUserController() user.Controller {
	wire.Build(userFeatureSet)
	return nil
}
