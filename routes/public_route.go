package routes

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/user"
)

func PublicRoute(routerGroup *gin.RouterGroup, userController user.Controller) {
	routerGroup.POST("/register", userController.Register)
	routerGroup.POST("/generate-otp", userController.GenerateOneTimePassword)
	routerGroup.POST("/verify-otp", userController.VerifyOneTimePassword)
}
