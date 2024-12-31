package note

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/note/dto"
)

type Service interface {
	HandleCreate(ginContext *gin.Context, createNoteDto *dto.CreateNoteDto)
	HandleUpdate(ginContext *gin.Context, updateNoteDto *dto.UpdateNoteDto)
	HandleDelete(ginContext *gin.Context, noteId *string)
}
