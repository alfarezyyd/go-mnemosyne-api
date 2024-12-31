package whatsapp

import "github.com/gin-gonic/gin"

type Controller interface {
	VerifyTokenWebhook(ginContext *gin.Context)
	Create(ginContext *gin.Context)
	ProcessWebhook(ginContext *gin.Context)
}
