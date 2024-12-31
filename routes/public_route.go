package routes

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/whatsapp"
)

func PublicRoute(routerGroup *gin.RouterGroup, whatsAppController whatsapp.Controller) {
	whatsAppGroup := routerGroup.Group("/whatsapp")
	whatsAppGroup.GET("/webhook", whatsAppController.VerifyTokenWebhook)
	whatsAppGroup.POST("/webhook", whatsAppController.ProcessWebhook)
}
