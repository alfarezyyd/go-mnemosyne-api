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
	allCategoryByUser := categoryHandler.categoryService.HandleGetAllByUser(ginContext)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Category fetch successfully", allCategoryByUser))
}

func (categoryHandler *Handler) Create(ginContext *gin.Context) {
	var categoryCreateDto dto.CreateCategoryDto
	err := ginContext.ShouldBindBodyWithJSON(&categoryCreateDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	categoryHandler.categoryService.HandleCreate(ginContext, &categoryCreateDto)
	ginContext.JSON(http.StatusCreated, helper.WriteSuccess("Category has been created", nil))
}

func (categoryHandler *Handler) Update(ginContext *gin.Context) {
	var updateCategoryDto dto.UpdateCategoryDto
	err := ginContext.ShouldBindBodyWithJSON(&updateCategoryDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	categoryHandler.categoryService.HandleUpdate(ginContext, &updateCategoryDto)
	ginContext.JSON(http.StatusCreated, helper.WriteSuccess("Category has been updated", nil))
}

func (categoryHandler *Handler) Delete(ginContext *gin.Context) {
	categoryId := ginContext.Param("id")
	categoryHandler.categoryService.HandleDelete(ginContext, categoryId)
	ginContext.JSON(http.StatusOK, helper.WriteSuccess("Category has been deleted", nil))
}
