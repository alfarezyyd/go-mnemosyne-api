package note

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/note/dto"
)

type Service interface {
	HandleCreate(ginContext *gin.Context, createNoteDto *dto.CreateNoteDto)
}
