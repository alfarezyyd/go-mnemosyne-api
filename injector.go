//go:build wireinject
// +build wireinject

package main

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go-mnemosyne-api/category"
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

var categoryFeatureSet = wire.NewSet(
	category.NewHandler,
	wire.Bind(new(category.Controller), new(*category.Handler)),
	category.NewService,
	wire.Bind(new(category.Service), new(*category.ServiceImpl)),
	category.NewRepository,
	wire.Bind(new(category.Repository), new(*category.RepositoryImpl)),
)

func InitializeUserController(dbConnection *gorm.DB,
	validatorInstance *validator.Validate,
	engTranslator ut.Translator,
	mailerService *config.MailerService,
	identityProvider *config.IdentityProvider,
	viperConfig *viper.Viper) user.Controller {
	wire.Build(userFeatureSet)
	return nil
}

func InitializeCategoryController(
	dbConnection *gorm.DB,
	validatorInstance *validator.Validate,
	engTranslator ut.Translator) category.Controller {
	wire.Build(categoryFeatureSet)
	return nil
}
