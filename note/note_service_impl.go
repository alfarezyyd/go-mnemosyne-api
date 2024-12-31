package note

import (
	"fmt"
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

func NewService(noteRepository Repository, dbConnection *gorm.DB, validationService *validator.Validate, engTranslator ut.Translator) *ServiceImpl {
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
	fmt.Println(userJwtClaim)
	err = noteService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		var noteModel model.Note
		var isCategoryExists bool
		err = gormTransaction.Where("email = ?", userJwtClaim.Email).First(&userModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = gormTransaction.Model(&model.Category{}).Select("COUNT(*) > 0").Where("id = ?", createNoteDto.CategoryId).Where("user_id = ?", userModel.ID).Find(&isCategoryExists).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapNoteDtoIntoNoteModel(createNoteDto, &noteModel)
		noteModel.UserID = userModel.ID
		err = gormTransaction.Create(&noteModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
}

func (noteService *ServiceImpl) HandleUpdate(ginContext *gin.Context, updateNoteDto *dto.UpdateNoteDto) {
	err := noteService.validationService.Struct(updateNoteDto)
	exception.ParseValidationError(err, noteService.engTranslator)
	userJwtClaim := ginContext.MustGet("claims").(*userDto.JwtClaimDto)
	noteId := ginContext.Param("id")
	err = noteService.validationService.Var(noteId, "required,gte=1")
	exception.ParseValidationError(err, noteService.engTranslator)
	err = noteService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		var existingNote model.Note
		var isCategoryExists bool
		err = gormTransaction.Where("email = ?", userJwtClaim.Email).First(&userModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = gormTransaction.Where("id = ?", noteId).First(&existingNote).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		if existingNote.CategoryId != updateNoteDto.CategoryId {
			err = gormTransaction.Model(&model.Category{}).Select("COUNT(*) > 0").Where("id = ?", updateNoteDto.CategoryId).Where("user_id = ?", userModel.ID).Find(&isCategoryExists).Error
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
			existingNote.CategoryId = updateNoteDto.CategoryId
		}
		mapper.MapNoteDtoIntoNoteModel(updateNoteDto, &existingNote)
		err = gormTransaction.Where("id = ?", existingNote.ID).Updates(&existingNote).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
}

func (noteService *ServiceImpl) HandleDelete(ginContext *gin.Context, noteId *string) {
	err := noteService.validationService.Var(noteId, "required,gte=1")
	exception.ParseValidationError(err, noteService.engTranslator)
	userJwtClaim := ginContext.MustGet("claims").(*userDto.JwtClaimDto)
	err = noteService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		err = gormTransaction.Where("email = ?", userJwtClaim.Email).First(&userModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = gormTransaction.
			Where("notes.id = ? AND user_id = ?", noteId, userModel.ID).
			Delete(&model.Note{}).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
}
