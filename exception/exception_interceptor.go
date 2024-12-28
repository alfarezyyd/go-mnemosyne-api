package exception

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Interceptor() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		defer func() {
			if occurredError := recover(); occurredError != nil {
				// Check if it's our custom
				var clientError *ClientError
				if errors.As(occurredError.(*ClientError), &clientError) {
					ginContext.AbortWithStatusJSON(clientError.StatusCode, gin.H{
						"message": clientError.Message,
					})
					return
				}

				// Unknown
				ginContext.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Internal server occurred",
				})
			}
		}()
		ginContext.Next()
	}
}
