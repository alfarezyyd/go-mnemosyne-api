package sharedNote

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	dbConnection         *gorm.DB
	sharedNoteRepository Repository
	validatorInstance    *validator.Validate
	engTranslator        *ut.Translator
}

func NewService(dbConnection *gorm.DB,
	sharedNoteRepository Repository,
	validatorInstance *validator.Validate,
	engTranslator *ut.Translator) *ServiceImpl {
	return &ServiceImpl{
		dbConnection:         dbConnection,
		sharedNoteRepository: sharedNoteRepository,
		validatorInstance:    validatorInstance,
		engTranslator:        engTranslator,
	}
}
