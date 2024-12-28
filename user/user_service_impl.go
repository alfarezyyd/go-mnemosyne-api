package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/user/dto"
	"gorm.io/gorm"
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
	validationError := exception.ParseValidationError(err, serviceImpl.engTranslator)
	fmt.Println(validationError)
}
