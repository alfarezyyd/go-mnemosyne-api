package whatsapp

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/whatsapp/dto"
	"net/http"
)

type Handler struct {
	whatsAppService Service
}

func NewHandler(whatsAppService Service) *Handler {
	return &Handler{whatsAppService: whatsAppService}
}

func (whatsAppHandler *Handler) VerifyTokenWebhook(ginContext *gin.Context) {
	whatsAppHandler.whatsAppService.HandleVerifyTokenWebhook(ginContext)
}

func (whatsAppHandler *Handler) ProcessWebhook(ginContext *gin.Context) {
	var payloadMessageDto dto.PayloadMessageDto
	err := ginContext.ShouldBindBodyWithJSON(&payloadMessageDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	whatsAppHandler.whatsAppService.HandleMessageWebhook(ginContext, &payloadMessageDto)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Success", nil))
}

func (whatsAppHandler *Handler) Create(ginContext *gin.Context) {
	whatsAppHandler.whatsAppService.HandleCreate(ginContext)
}
