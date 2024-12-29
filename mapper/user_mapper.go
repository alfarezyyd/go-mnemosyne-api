package mapper

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-viper/mapstructure/v2"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/model"
	"go-mnemosyne-api/user/dto"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func MapUserDtoIntoUserModel[T *dto.CreateUserDto](userTransferObject T) *model.User {
	var userModel model.User
	err := mapstructure.Decode(userTransferObject, &userModel)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), 14)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	userModel.Password = string(hashedPassword)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	return &userModel
}

func MapJwtClaimIntoUserClaim(jwtClaim jwt.MapClaims) (*dto.JwtClaimDto, error) {
	var userClaim dto.JwtClaimDto
	err := mapstructure.Decode(jwtClaim, &userClaim)
	if err != nil {
		return nil, err
	}
	return &userClaim, nil
}
