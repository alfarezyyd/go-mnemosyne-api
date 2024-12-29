//go:build wireinject
// +build wireinject

package main

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"go-mnemosyne-api/config"
	"go-mnemosyne-api/user"
	"gorm.io/gorm"
)

var userFeatureSet = wire.NewSet(
	user.NewHandler,
	wire.Bind(new(user.Controller), new(*user.Handler)),
	user.NewService,
	wire.Bind(new(user.Service), new(*user.ServiceImpl)),
	user.NewRepository,
	wire.Bind(new(user.Repository), new(*user.RepositoryImpl)),
)

func InitializeUserController(dbConnection *gorm.DB,
	validatorInstance *validator.Validate,
	engTranslator ut.Translator,
	mailerService *config.MailerService) user.Controller {
	wire.Build(userFeatureSet)
	return nil
}
