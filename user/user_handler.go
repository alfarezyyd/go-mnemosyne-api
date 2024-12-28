package user

import "github.com/gin-gonic/gin"

type Handler struct {
	userService Service
}

func NewHandler(userService Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (handler *Handler) Register(ginContext *gin.Context) {}
