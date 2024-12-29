package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go-mnemosyne-api/config"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/mapper"
	"go-mnemosyne-api/model"
	"go-mnemosyne-api/user/dto"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

type ServiceImpl struct {
	userRepository    Repository
	dbConnection      *gorm.DB
	validatorInstance *validator.Validate
	engTranslator     ut.Translator
	mailerService     *config.MailerService
}

func NewService(userRepository Repository, dbConnection *gorm.DB, validatorInstance *validator.Validate, engTranslator ut.Translator, mailerService *config.MailerService) *ServiceImpl {
	return &ServiceImpl{
		userRepository:    userRepository,
		dbConnection:      dbConnection,
		validatorInstance: validatorInstance,
		engTranslator:     engTranslator,
		mailerService:     mailerService,
	}
}

func (serviceImpl *ServiceImpl) HandleRegister(ginContext *gin.Context, createUserDto *dto.CreateUserDto) {
	err := serviceImpl.validatorInstance.Struct(createUserDto)
	exception.ParseValidationError(err, serviceImpl.engTranslator)

	err = serviceImpl.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		userModel := mapper.MapUserDtoIntoUserModel(createUserDto)
		err = gormTransaction.Create(userModel).Error
		helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (serviceImpl *ServiceImpl) HandleGenerateOneTimePassword(ginContext *gin.Context, generateOneTimePassDto *dto.GenerateOtpDto) {
	err := serviceImpl.validatorInstance.Struct(generateOneTimePassDto)
	exception.ParseValidationError(err, serviceImpl.engTranslator)
	err = serviceImpl.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		var oneTimePasswordToken model.OneTimePasswordToken
		err = gormTransaction.Where("email = ?", generateOneTimePassDto.Email).First(&userModel).Error
		generatedOneTimePasswordToken := helper.GenerateOneTimePasswordToken()
		hashedGeneratedOneTimePasswordToken, err := bcrypt.GenerateFromPassword([]byte(generatedOneTimePasswordToken), 10)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		oneTimePasswordToken.UserId = userModel.ID
		oneTimePasswordToken.HashedToken = string(hashedGeneratedOneTimePasswordToken)
		oneTimePasswordToken.ExpiresAt = time.Now().Add(15 * time.Minute)
		err = gormTransaction.Create(&oneTimePasswordToken).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		emailPayload := config.EmailPayload{
			Title:     "OTP Sent",
			Recipient: generateOneTimePassDto.Email,
			Body:      fmt.Sprintf("One Time Password %s", generatedOneTimePasswordToken),
			Sender:    serviceImpl.mailerService.ViperConfig.GetString(""),
		}

		projectRoot, _ := os.Getwd() // Mendapatkan root path proyek
		templateFile := fmt.Sprintf("%s/public/static/email_template.html", projectRoot)
		err = serviceImpl.mailerService.SendEmail(
			generateOneTimePassDto.Email,
			"OTP Send",
			templateFile,
			emailPayload)
		helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
		return nil
	})
}
