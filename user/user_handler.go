package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/user/dto"
	"net/http"
)

type Handler struct {
	userService       Service
	validatorInstance *validator.Validate
}

func NewHandler(userService Service, validatorInstance *validator.Validate) *Handler {
	return &Handler{
		userService:       userService,
		validatorInstance: validatorInstance,
	}
}

func (userHandler *Handler) Register(ginContext *gin.Context) {
	var createUserDto dto.CreateUserDto
	err := ginContext.ShouldBindBodyWithJSON(&createUserDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	userHandler.userService.HandleRegister(ginContext, &createUserDto)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("User created successfully", nil))
}

func (userHandler *Handler) GenerateOneTimePassword(ginContext *gin.Context) {
	var generateOneTimePassDto dto.GenerateOtpDto
	err := ginContext.ShouldBindBodyWithJSON(&generateOneTimePassDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	userHandler.userService.HandleGenerateOneTimePassword(ginContext, &generateOneTimePassDto)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("OTP generated successfully", nil))
}

func (userHandler *Handler) VerifyOneTimePassword(ginContext *gin.Context) {
	var VerifyOneTimePassDto dto.VerifyOtpDto
	err := ginContext.ShouldBindBodyWithJSON(&VerifyOneTimePassDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	userHandler.userService.HandleVerifyOneTimePassword(ginContext, &VerifyOneTimePassDto)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("OTP verified successfully", nil))
}
