package user

import "github.com/gin-gonic/gin"

type Controller interface {
	Register(ginContext *gin.Context)
}
