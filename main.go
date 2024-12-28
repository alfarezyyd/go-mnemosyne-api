package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-mnemosyne-api/config"
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
	databaseInstance.GetDatabaseConnection()

	// Gin Initialization
	ginEngine := gin.Default()
	ginEngine.Use(gin.Recovery())
	ginError := ginEngine.Run()
	if ginError != nil {
		panic(ginError)
	}
}
