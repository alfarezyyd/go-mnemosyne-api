package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/note/dto"
	"net/http"
)

type Handler struct {
	noteService Service
}

func NewHandler(noteService Service) *Handler {
	return &Handler{
		noteService: noteService,
	}
}

func (noteHandler *Handler) Create(ginContext *gin.Context) {
	var createNoteDto dto.CreateNoteDto
	err := ginContext.ShouldBindBodyWithJSON(&createNoteDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	fmt.Println(createNoteDto)
	noteHandler.noteService.HandleCreate(ginContext, &createNoteDto)
	ginContext.JSON(http.StatusCreated, helper.WriteSuccess("Success", nil))
}
