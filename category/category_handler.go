package category

import (
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/category/dto"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"net/http"
)

type Handler struct {
	categoryService Service
}

func NewHandler(categoryService Service) *Handler {
	return &Handler{
		categoryService: categoryService,
	}
}

func (categoryHandler *Handler) GetAllByUser(ginContext *gin.Context) {

}

func (categoryHandler *Handler) Create(ginContext *gin.Context) {
	var categoryCreateDto dto.CreateCategoryDto
	err := ginContext.ShouldBindBodyWithJSON(&categoryCreateDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	categoryHandler.categoryService.HandleCreate(ginContext, &categoryCreateDto)
}

func (categoryHandler *Handler) Update(ginContext *gin.Context) {}

func (categoryHandler *Handler) Delete(ginContext *gin.Context) {}
