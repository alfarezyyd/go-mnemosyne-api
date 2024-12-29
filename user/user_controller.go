package user

import "github.com/gin-gonic/gin"

type Controller interface {
	Register(ginContext *gin.Context)
	GenerateOneTimePassword(ginContext *gin.Context)
	VerifyOneTimePassword(ginContext *gin.Context)
	LoginWithGoogle(ginContext *gin.Context)
	GoogleProviderCallback(ginContext *gin.Context)
}
