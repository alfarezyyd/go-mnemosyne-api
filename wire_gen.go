// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go-mnemosyne-api/category"
	"go-mnemosyne-api/config"
	"go-mnemosyne-api/discord"
	"go-mnemosyne-api/note"
	"go-mnemosyne-api/shared_note"
	"go-mnemosyne-api/user"
	"go-mnemosyne-api/whatsapp"
	"gorm.io/gorm"
)

// Injectors from injector.go:

func InitializeUserController(dbConnection *gorm.DB, validatorInstance *validator.Validate, engTranslator ut.Translator, mailerService *config.MailerService, identityProvider *config.IdentityProvider, viperConfig *viper.Viper) user.Controller {
	repositoryImpl := user.NewRepository()
	serviceImpl := user.NewService(repositoryImpl, dbConnection, validatorInstance, engTranslator, mailerService, identityProvider, viperConfig)
	handler := user.NewHandler(serviceImpl, validatorInstance)
	return handler
}

func InitializeWhatsAppController(dbConnection *gorm.DB, validatorInstance *validator.Validate, engTranslator ut.Translator, viperConfig *viper.Viper, vertexClient *config.VertexClient, cloudStorageClient *config.GoogleCloudStorage) whatsapp.Controller {
	repositoryImpl := whatsapp.NewRepository()
	noteRepositoryImpl := note.NewRepository()
	serviceImpl := note.NewService(noteRepositoryImpl, dbConnection, validatorInstance, engTranslator)
	whatsappServiceImpl := whatsapp.NewService(repositoryImpl, dbConnection, viperConfig, engTranslator, vertexClient, serviceImpl, cloudStorageClient)
	handler := whatsapp.NewHandler(whatsappServiceImpl)
	return handler
}

func InitializeCategoryController(dbConnection *gorm.DB, validatorInstance *validator.Validate, engTranslator ut.Translator) category.Controller {
	repositoryImpl := category.NewRepository()
	serviceImpl := category.NewService(repositoryImpl, dbConnection, validatorInstance, engTranslator)
	handler := category.NewHandler(serviceImpl)
	return handler
}

func InitializeNoteController(dbConnection *gorm.DB, validatorInstance *validator.Validate, engTranslator ut.Translator) note.Controller {
	repositoryImpl := note.NewRepository()
	serviceImpl := note.NewService(repositoryImpl, dbConnection, validatorInstance, engTranslator)
	handler := note.NewHandler(serviceImpl)
	return handler
}

func InitializeSharedNoteController(dbConnection *gorm.DB, validatorInstance *validator.Validate, engTranslator ut.Translator) sharedNote.Controller {
	handler := sharedNote.NewHandler()
	return handler
}

func InitializeDiscordController(dbConnection *gorm.DB, vertexClient *config.VertexClient, validatorInstance *validator.Validate, engTranslator ut.Translator, cloudStorageClient *config.GoogleCloudStorage, viperConfig *viper.Viper) discord.Controller {
	repositoryImpl := discord.NewRepository()
	serviceImpl := discord.NewService(repositoryImpl)
	noteRepositoryImpl := note.NewRepository()
	noteServiceImpl := note.NewService(noteRepositoryImpl, dbConnection, validatorInstance, engTranslator)
	handler := discord.NewHandler(serviceImpl, vertexClient, noteServiceImpl, dbConnection, cloudStorageClient, viperConfig)
	return handler
}

// injector.go:

var userFeatureSet = wire.NewSet(user.NewHandler, wire.Bind(new(user.Controller), new(*user.Handler)), user.NewService, wire.Bind(new(user.Service), new(*user.ServiceImpl)), user.NewRepository, wire.Bind(new(user.Repository), new(*user.RepositoryImpl)))

var categoryFeatureSet = wire.NewSet(category.NewHandler, wire.Bind(new(category.Controller), new(*category.Handler)), category.NewService, wire.Bind(new(category.Service), new(*category.ServiceImpl)), category.NewRepository, wire.Bind(new(category.Repository), new(*category.RepositoryImpl)))

var noteFeatureSet = wire.NewSet(note.NewHandler, wire.Bind(new(note.Controller), new(*note.Handler)), note.NewService, wire.Bind(new(note.Service), new(*note.ServiceImpl)), note.NewRepository, wire.Bind(new(note.Repository), new(*note.RepositoryImpl)))

var sharedNoteFeatureSet = wire.NewSet(sharedNote.NewHandler, wire.Bind(new(sharedNote.Controller), new(*sharedNote.Handler)), sharedNote.NewService, wire.Bind(new(sharedNote.Service), new(*sharedNote.ServiceImpl)), sharedNote.NewRepository, wire.Bind(new(sharedNote.Repository), new(*sharedNote.RepositoryImpl)))

var whatsAppFeatureSet = wire.NewSet(whatsapp.NewHandler, wire.Bind(new(whatsapp.Controller), new(*whatsapp.Handler)), whatsapp.NewService, wire.Bind(new(whatsapp.Service), new(*whatsapp.ServiceImpl)), whatsapp.NewRepository, wire.Bind(new(whatsapp.Repository), new(*whatsapp.RepositoryImpl)))

var discordFeatureSet = wire.NewSet(discord.NewHandler, wire.Bind(new(discord.Controller), new(*discord.Handler)), discord.NewService, wire.Bind(new(discord.Service), new(*discord.ServiceImpl)), discord.NewRepository, wire.Bind(new(discord.Repository), new(*discord.RepositoryImpl)))
