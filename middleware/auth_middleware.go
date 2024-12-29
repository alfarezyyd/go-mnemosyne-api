package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/mapper"
	"net/http"
	"strings"
)

func AuthMiddleware(viperConfig *viper.Viper) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		tokenString := ginContext.GetHeader("Authorization")
		trimmedTokenString := strings.Replace(tokenString, "Bearer ", "", 1)
		// Parse the token
		token, err := jwt.Parse(trimmedTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(viperConfig.GetString("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ginContext.Abort() // Stop further processing if unauthorized
			return
		}

		// Set the token claims to the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userJwtClaim, err := mapper.MapJwtClaimIntoUserClaim(claims)
			helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
			ginContext.Set("claims", userJwtClaim)
		} else {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ginContext.Abort()
			return
		}

		ginContext.Next() // Proceed to the next handler if authorized
	}
}
