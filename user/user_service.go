package user

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/user/dto"
)

type Service interface {
	HandleRegister(ginContext *gin.Context, createUserDto *dto.CreateUserDto)
	HandleGenerateOneTimePassword(ginContext *gin.Context, generateOtpDto *dto.GenerateOtpDto)
	HandleVerifyOneTimePassword(ginContext *gin.Context, verifyOtpDto *dto.VerifyOtpDto)
	HandleGoogleAuthentication(ginContext *gin.Context)
	HandleGoogleCallback(ginContext *gin.Context)
	HandleLogin(ginContext *gin.Context, loginUserDto *dto.LoginUserDto) string
}
