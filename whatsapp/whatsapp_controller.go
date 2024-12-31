package whatsapp

import "github.com/gin-gonic/gin"

type Controller interface {
	Create(ginContext *gin.Context)
}
