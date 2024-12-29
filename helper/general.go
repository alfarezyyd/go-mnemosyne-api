package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-mnemosyne-api/exception"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"strconv"
)

func CheckErrorOperation(indicatedError error, clientError *exception.ClientError) bool {
	if indicatedError != nil {
		panic(clientError)
		return true
	}
	return false
}

func TransactionOperation(runningTransaction *gorm.DB, ginContext *gin.Context) {
	occurredError := recover()
	fmt.Println(occurredError)
	if occurredError != nil {
		fmt.Println(occurredError)
		runningTransaction.Rollback()
		ginContext.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": runningTransaction.Error.Error()})
	} else {
		runningTransaction.Commit()
	}
}

func GenerateOneTimePasswordToken() string {
	num := rand.Intn(9000) + 1000
	return strconv.Itoa(num)
}
