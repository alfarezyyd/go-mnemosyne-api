package mapper

import (
	"github.com/go-viper/mapstructure/v2"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/model"
	"go-mnemosyne-api/note/dto"
	"net/http"
)

func MapNoteDtoIntoNoteModel[T *dto.CreateNoteDto](noteDto T, noteModel *model.Note) {
	err := mapstructure.Decode(noteDto, noteModel)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
}
