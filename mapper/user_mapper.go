package mapper

import (
	"github.com/go-viper/mapstructure/v2"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/model"
	"go-mnemosyne-api/user/dto"
	"net/http"
)

func MapUserDtoIntoUserModel[T *dto.CreateUserDto](userTransferObject T) *model.User {
	var userModel model.User
	err := mapstructure.Decode(userTransferObject, &userModel)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	return &userModel
}
