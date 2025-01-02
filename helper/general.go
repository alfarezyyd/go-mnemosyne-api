package helper

import (
	"fmt"
	"go-mnemosyne-api/exception"
	"math/rand"
	"reflect"
	"strconv"
)

func CheckErrorOperation(indicatedError error, clientError *exception.ClientError) bool {
	if indicatedError != nil {
		panic(clientError)
		return true
	}
	return false
}

func ParseNullableValue(value interface{}) interface{} {
	v := reflect.ValueOf(value)

	// Handle pointer types
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check for sql.Null* types
	if v.Kind() == reflect.Struct {
		validField := v.FieldByName("Valid")
		valueField := v.FieldByName("String") // Default to String; change for other nullable types

		if validField.IsValid() && validField.Kind() == reflect.Bool {
			if validField.Bool() {
				return valueField.Interface()
			}
			return nil
		}
	}

	// If not a nullable type, return the value as is
	return value
}

func GenerateOneTimePasswordToken() string {
	num := rand.Intn(9000) + 1000
	return strconv.Itoa(num)
}

func LogError(err error) {
	if err != nil {
		fmt.Println("error occurred", err)
	}
}
