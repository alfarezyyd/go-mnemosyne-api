package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-mnemosyne-api/config"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/middleware"
	"go-mnemosyne-api/routes"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	viperConfig := viper.New()
	viperConfig.SetConfigFile(".env")
	viperConfig.AddConfigPath(".")
	viperConfig.AutomaticEnv()
	viperConfig.ReadInConfig()

	// Database Initialization
	databaseCredentials := &config.DatabaseCredentials{
		DatabaseHost:     viperConfig.GetString("DATABASE_HOST"),
		DatabasePort:     viperConfig.GetString("DATABASE_PORT"),
		DatabaseName:     viperConfig.GetString("DATABASE_NAME"),
		DatabasePassword: viperConfig.GetString("DATABASE_PASSWORD"),
		DatabaseUsername: viperConfig.GetString("DATABASE_USERNAME"),
	}

	databaseInstance := config.NewDatabaseConnection(databaseCredentials)
	databaseConnection := databaseInstance.GetDatabaseConnection()

	// VertexAI Client
	vertexInstance := config.NewVertexClient(viperConfig)
	err := vertexInstance.InitializeVertexClient()
	if err != nil {
		panic(err)
	}

	storage, err := config.InitializeGoogleCloudStorage()
	if err != nil {
		panic(err)
	}
	// Validator
	validatorInstance, engTranslator := config.InitializeValidator()
	discordClient := config.NewDiscordClient(viperConfig)
	discordSession, err := discordClient.InitializeDiscordConnection()
	if err != nil {
		panic(err)
	}
	mailerService := config.NewMailerService(viperConfig)
	// Gin Initialization
	ginEngine := gin.Default()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(exception.Interceptor())

	identityProvider := config.NewIdentityProvider(viperConfig)
	identityProvider.InitializeGoogleProviderConfig()

	userController := InitializeUserController(databaseConnection, validatorInstance, engTranslator, mailerService, identityProvider, viperConfig)
	categoryController := InitializeCategoryController(databaseConnection, validatorInstance, engTranslator)
	noteController := InitializeNoteController(databaseConnection, validatorInstance, engTranslator)
	whatsAppController := InitializeWhatsAppController(databaseConnection, validatorInstance, engTranslator, viperConfig, vertexInstance, storage)
	discordController := InitializeDiscordController(databaseConnection, vertexInstance, validatorInstance, engTranslator, storage, viperConfig)
	authRouterGroup := ginEngine.Group("/authentication")
	routes.AuthenticationRoute(authRouterGroup, userController)
	publicRouterGroup := ginEngine.Group("/public")
	routes.PublicRoute(publicRouterGroup, whatsAppController)
	routes.DiscordRoutes(discordSession, discordController)
	apiRouterGroup := ginEngine.Group("/api")
	apiRouterGroup.Use(middleware.AuthMiddleware(viperConfig))
	routes.UserRoute(apiRouterGroup, categoryController, noteController)
	// Route
	ginError := ginEngine.Run()
	if ginError != nil {
		panic(ginError)
	}
}
