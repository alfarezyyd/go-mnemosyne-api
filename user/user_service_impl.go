package user

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go-mnemosyne-api/config"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/mapper"
	"go-mnemosyne-api/model"
	"go-mnemosyne-api/user/dto"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
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
	identityProvider  *config.IdentityProvider
	viperConfig       *viper.Viper
}

func NewService(
	userRepository Repository,
	dbConnection *gorm.DB,
	validatorInstance *validator.Validate,
	engTranslator ut.Translator,
	mailerService *config.MailerService,
	identityProvider *config.IdentityProvider,
	viperConfig *viper.Viper) *ServiceImpl {
	return &ServiceImpl{
		userRepository:    userRepository,
		dbConnection:      dbConnection,
		validatorInstance: validatorInstance,
		engTranslator:     engTranslator,
		mailerService:     mailerService,
		identityProvider:  identityProvider,
		viperConfig:       viperConfig,
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

func (serviceImpl *ServiceImpl) HandleVerifyOneTimePassword(ginContext *gin.Context, verifyOtpDto *dto.VerifyOtpDto) {
	err := serviceImpl.validatorInstance.Struct(verifyOtpDto)
	exception.ParseValidationError(err, serviceImpl.engTranslator)
	err = serviceImpl.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		var oneTimePasswordToken model.OneTimePasswordToken
		err = gormTransaction.Where("email = ?", verifyOtpDto.Email).First(&userModel).Error
		gormTransaction.Where("user_id = ?", userModel.ID).Order("expires_at desc").First(&oneTimePasswordToken)
		if !(time.Now().Before(oneTimePasswordToken.ExpiresAt)) {
			exception.ThrowClientError(exception.NewClientError(http.StatusBadRequest, "OTP has expired"))
		}
		fmt.Println(oneTimePasswordToken.HashedToken, verifyOtpDto.OneTimePasswordToken)
		err := bcrypt.CompareHashAndPassword([]byte(oneTimePasswordToken.HashedToken), []byte(verifyOtpDto.OneTimePasswordToken))
		helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
		return nil
	})
}

func (serviceImpl *ServiceImpl) HandleGoogleAuthentication(ginContext *gin.Context) {
	authEndpoint := serviceImpl.identityProvider.GoogleProviderConfig.AuthCodeURL("randomstate")
	ginContext.Redirect(http.StatusSeeOther, authEndpoint)
	ginContext.JSON(http.StatusOK, authEndpoint)
}

func (serviceImpl *ServiceImpl) HandleGoogleCallback(ginContext *gin.Context) {
	state := ginContext.Query("state")
	if state != "randomstate" {
		ginContext.JSON(http.StatusBadRequest, exception.ErrBadRequest)
		return
	}

	queryCode := ginContext.Query("code")

	googleProviderConfig := serviceImpl.identityProvider.GoogleProviderConfig

	token, err := googleProviderConfig.Exchange(context.Background(), queryCode)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, exception.ErrBadRequest)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, exception.ErrBadRequest)
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, exception.ErrBadRequest)
	}
	ginContext.JSON(http.StatusOK, string(userData))
}

func (serviceImpl *ServiceImpl) HandleLogin(ginContext *gin.Context, loginUserDto *dto.LoginUserDto) string {
	err := serviceImpl.validatorInstance.Struct(loginUserDto)
	exception.ParseValidationError(err, serviceImpl.engTranslator)
	var tokenString string
	err = serviceImpl.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		err = gormTransaction.Where("email = ?", loginUserDto.UserIdentifier).Or("phone_number = ?", loginUserDto.UserIdentifier).First(&userModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(loginUserDto.Password))
		helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, "User credentials invalid"))
		tokenInstance := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email":        userModel.Email,
			"phone_number": userModel.PhoneNumber,
			"exp":          time.Now().Add(time.Hour * 72).Unix(),
		})
		tokenString, err = tokenInstance.SignedString([]byte(serviceImpl.viperConfig.GetString("JWT_SECRET")))
		helper.CheckErrorOperation(err, exception.NewClientError(http.StatusInternalServerError, exception.ErrInternalServerError))
		return nil
	})
	return tokenString
}
