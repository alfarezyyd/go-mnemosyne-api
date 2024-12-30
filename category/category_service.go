package category

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/category/dto"
)

type Service interface {
	HandleCreate(ginContext *gin.Context, categoryCreateDto *dto.CreateCategoryDto)
	HandleUpdate(ginContext *gin.Context, updateCategoryDto *dto.UpdateCategoryDto)
	HandleDelete(ginContext *gin.Context, categoryId string)
}
