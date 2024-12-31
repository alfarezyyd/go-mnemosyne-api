package category

import (
	"fmt"
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

func (serviceImpl *ServiceImpl) HandleGetAllByUser(ginContext *gin.Context) []model.Category {
	var allCategoryById []model.Category
	userJwtClaim := ginContext.MustGet("claims").(*userDto.JwtClaimDto)
	err := serviceImpl.dbConnection.
		Joins("JOIN users ON users.id = categories.user_id").
		Select("categories.*").
		Where("users.email = ?", userJwtClaim.Email).
		Find(&allCategoryById).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	fmt.Println(allCategoryById)
	return allCategoryById
}
func (serviceImpl *ServiceImpl) HandleCreate(ginContext *gin.Context, categoryCreateDto *dto.CreateCategoryDto) {
	err := serviceImpl.validationInstance.Struct(categoryCreateDto)
	exception.ParseValidationError(err, serviceImpl.engTranslator)
	userJwtClaim := ginContext.MustGet("claims").(*userDto.JwtClaimDto)

	err = serviceImpl.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		var categoryModel model.Category
		err = gormTransaction.Where("email = ?", *userJwtClaim.Email).First(&userModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapCategoryDtoIntoCategoryModel(&categoryModel, categoryCreateDto)
		categoryModel.UserID = userModel.ID
		err = gormTransaction.Create(&categoryModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (serviceImpl *ServiceImpl) HandleUpdate(ginContext *gin.Context, updateCategoryDto *dto.UpdateCategoryDto) {
	err := serviceImpl.validationInstance.Struct(updateCategoryDto)
	exception.ParseValidationError(err, serviceImpl.engTranslator)
	queryParam := ginContext.Param("id")
	err = serviceImpl.validationInstance.Var(queryParam, "required,gte=1")
	exception.ParseValidationError(err, serviceImpl.engTranslator)
	userJwtClaim := ginContext.MustGet("claims").(*userDto.JwtClaimDto)
	err = serviceImpl.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var categoryModel model.Category
		err = gormTransaction.
			Joins("JOIN users ON users.id = categories.user_id").
			Where("categories.id = ?", queryParam).
			Where("users.email = ?", userJwtClaim.Email).
			First(&categoryModel).
			Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapCategoryDtoIntoCategoryModel(&categoryModel, updateCategoryDto)
		err = gormTransaction.Where("id = ?", queryParam).Updates(updateCategoryDto).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
}

func (serviceImpl *ServiceImpl) HandleDelete(ginContext *gin.Context, userId string) {
	err := serviceImpl.validationInstance.Var(userId, "required,number,gte=1")
	exception.ParseValidationError(err, serviceImpl.engTranslator)
	userJwtClaim := ginContext.MustGet("claims").(*userDto.JwtClaimDto)
	err = serviceImpl.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var categoryModel model.Category
		err = gormTransaction.
			Joins("JOIN users ON users.id = categories.user_id").
			Where("categories.id = ?", userId).
			Where("users.email = ?", userJwtClaim.Email).
			Delete(&categoryModel).
			Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
}
