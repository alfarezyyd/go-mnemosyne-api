package whatsapp

import "github.com/gin-gonic/gin"

type Handler struct {
	whatsAppService Service
}

func NewHandler(whatsAppService Service) *Handler {
	return &Handler{whatsAppService: whatsAppService}
}

func (serviceImpl *ServiceImpl) Create(ginContext *gin.Context) {

}
