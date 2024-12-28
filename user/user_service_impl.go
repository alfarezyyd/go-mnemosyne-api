package user

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/mapper"
	"go-mnemosyne-api/user/dto"
	"gorm.io/gorm"
	"net/http"
)

type ServiceImpl struct {
	userRepository    Repository
	dbConnection      *gorm.DB
	validatorInstance *validator.Validate
	engTranslator     ut.Translator
}

func NewService(userRepository Repository, dbConnection *gorm.DB, validatorInstance *validator.Validate, engTranslator ut.Translator) *ServiceImpl {
	return &ServiceImpl{
		userRepository:    userRepository,
		dbConnection:      dbConnection,
		validatorInstance: validatorInstance,
		engTranslator:     engTranslator,
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
