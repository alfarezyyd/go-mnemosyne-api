package routes

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/user"
)

func AuthenticationRoute(routerGroup *gin.RouterGroup, userController user.Controller) {
	routerGroup.GET("/google", userController.LoginWithGoogle)
	routerGroup.GET("/google/callback", userController.GoogleProviderCallback)

	routerGroup.POST("/login", userController.Login)
	routerGroup.POST("/register", userController.Register)
	routerGroup.POST("/generate-otp", userController.GenerateOneTimePassword)
	routerGroup.POST("/verify-otp", userController.VerifyOneTimePassword)
}
