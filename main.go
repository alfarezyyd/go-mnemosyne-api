package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-mnemosyne-api/config"
	"go-mnemosyne-api/exception"
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

	// Validator
	validatorInstance, engTranslator := config.InitializeValidator()

	mailerService := config.NewMailerService(viperConfig)
	// Gin Initialization
	ginEngine := gin.Default()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(exception.Interceptor())

	identityProvider := config.NewIdentityProvider(viperConfig)
	identityProvider.InitializeGoogleProviderConfig()

	userController := InitializeUserController(databaseConnection, validatorInstance, engTranslator, mailerService, identityProvider, viperConfig)
	authRouterGroup := ginEngine.Group("/authentication")
	routes.AuthenticationRoute(authRouterGroup, userController)

	apiRouterGroup := ginEngine.Group("/api")
	publicRouterGroup := apiRouterGroup.Group("/public")
	routes.PublicRoute(publicRouterGroup, userController)
	routes.UserRoute(apiRouterGroup, userController)
	// Route
	ginError := ginEngine.Run()
	if ginError != nil {
		panic(ginError)
	}
}
