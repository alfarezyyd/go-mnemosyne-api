package category

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go-mnemosyne-api/category/dto"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/mapper"
	"go-mnemosyne-api/model"
	userDto "go-mnemosyne-api/user/dto"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	categoryRepository Repository
	dbConnection       *gorm.DB
	validationInstance *validator.Validate
	engTranslator      ut.Translator
}

func NewService(
	categoryRepository Repository,
	dbConnection *gorm.DB,
	validatorInstance *validator.Validate,
	engTranslator ut.Translator) *ServiceImpl {
	return &ServiceImpl{
		categoryRepository: categoryRepository,
		dbConnection:       dbConnection,
		validationInstance: validatorInstance,
		engTranslator:      engTranslator,
	}
}

func (serviceImpl *ServiceImpl) HandleCreate(ginContext *gin.Context, categoryCreateDto *dto.CreateCategoryDto) {
	err := serviceImpl.validationInstance.Struct(categoryCreateDto)
	exception.ParseValidationError(err, serviceImpl.engTranslator)
	userJwtClaim := ginContext.MustGet("claims").(*userDto.JwtClaimDto)

	err = serviceImpl.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		err = gormTransaction.Where("email = ?", *userJwtClaim.Email).First(&userModel).Error

		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		categoryModel := mapper.MapCategoryDtoIntoCategoryModel(categoryCreateDto)
		categoryModel.UserID = userModel.ID
		err = gormTransaction.Create(&categoryModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
