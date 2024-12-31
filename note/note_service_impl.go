package note

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/mapper"
	"go-mnemosyne-api/model"
	"go-mnemosyne-api/note/dto"
	userDto "go-mnemosyne-api/user/dto"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	noteRepository    Repository
	dbConnection      *gorm.DB
	validationService *validator.Validate
	engTranslator     ut.Translator
}

func NewServiceImpl(noteRepository Repository, dbConnection *gorm.DB, validationService *validator.Validate, engTranslator ut.Translator) *ServiceImpl {
	return &ServiceImpl{
		noteRepository:    noteRepository,
		dbConnection:      dbConnection,
		validationService: validationService,
		engTranslator:     engTranslator,
	}
}

func (noteService *ServiceImpl) HandleCreate(ginContext *gin.Context, createNoteDto *dto.CreateNoteDto) {
	err := noteService.validationService.Struct(createNoteDto)
	exception.ParseValidationError(err, noteService.engTranslator)
	userJwtClaim := ginContext.MustGet("claims").(*userDto.JwtClaimDto)
	err = noteService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		var noteModel model.Note
		var isCategoryExists bool
		err = gormTransaction.Where("email = ?", userJwtClaim.Email).First(&userModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = gormTransaction.Model(&model.Category{}).Select("COUNT(*) > 0").Where("id = ?", createNoteDto.CategoryId).Find(&isCategoryExists).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapNoteDtoIntoNoteModel(createNoteDto, &noteModel)
		err = gormTransaction.Create(&noteModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
}
