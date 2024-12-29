package user

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/user/dto"
)

type Service interface {
	HandleRegister(ginContext *gin.Context, createUserDto *dto.CreateUserDto)
	HandleGenerateOneTimePassword(ginContext *gin.Context, generateOtpDto *dto.GenerateOtpDto)
}
