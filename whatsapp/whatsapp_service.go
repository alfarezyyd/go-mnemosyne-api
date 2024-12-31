package whatsapp

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/whatsapp/dto"
)

type Service interface {
	HandleCreate(ginContext *gin.Context)
	HandleVerifyTokenWebhook(ginContext *gin.Context)
	HandleProcessWebhook(ginContext *gin.Context, payloadMessageDto *dto.PayloadMessageDto)
}
