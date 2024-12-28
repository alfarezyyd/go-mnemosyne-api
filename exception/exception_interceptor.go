package exception

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/web"
	"net/http"
)

func Interceptor() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		defer func() {
			if occurredError := recover(); occurredError != nil {
				// Check if it's our custom
				var clientError *ClientError
				if errors.As(occurredError.(*ClientError), &clientError) {
					ginContext.AbortWithStatusJSON(clientError.StatusCode, web.NewResponseContract(
						false, clientError.Message, nil, clientError.Trace))
					return
				}

				// Unknown
				ginContext.AbortWithStatusJSON(http.StatusInternalServerError, web.NewResponseContract(
					false, "Internal server error", nil, nil))
			}
		}()
		ginContext.Next()
	}
}
