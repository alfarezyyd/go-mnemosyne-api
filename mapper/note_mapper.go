package mapper

import (
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/model"
	"go-mnemosyne-api/note/dto"
	"net/http"
)

func MapNoteDtoIntoNoteModel[T *dto.CreateNoteDto | *dto.UpdateNoteDto](noteDto T, noteModel *model.Note) {
	err := mapstructure.Decode(noteDto, noteModel)
	fmt.Println(err)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
}

func MapAllNoteIntoString(allNote []model.Note) string {
	var parsedString string
	for i, note := range allNote {
		parsedString += fmt.Sprintf("%d. %v \n", i+1, note.Title)
		if note.Content != "" {
			parsedString += fmt.Sprintf("Content: %v\n", note.Content)
		}
		if note.DueDate != nil {
			parsedString += fmt.Sprintf("Due Date: %v\n", *(note.DueDate))
		}
	}
	return parsedString
}
