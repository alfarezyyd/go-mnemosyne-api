package category

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/category/dto"
)

type Service interface {
	HandleCreate(ginContext *gin.Context, categoryCreateDto *dto.CreateCategoryDto)
}
