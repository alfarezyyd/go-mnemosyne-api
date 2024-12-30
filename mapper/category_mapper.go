package mapper

import (
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"go-mnemosyne-api/category/dto"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/model"
	"net/http"
)

func MapCategoryDtoIntoCategoryModel[T *dto.CreateCategoryDto | *dto.UpdateCategoryDto](categoryModel *model.Category, categoryDto T) {
	err := mapstructure.Decode(categoryDto, &categoryModel)
	fmt.Println(err)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
}
