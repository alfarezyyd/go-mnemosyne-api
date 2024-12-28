package exception

import "fmt"

type ClientError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (clientError *ClientError) Error() string {
	return fmt.Sprintf("Error %d: %s", clientError.StatusCode, clientError.Message)
}

func NewClientError(statusCode int, message string) *ClientError {
	return &ClientError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func ThrowClientError(clientError *ClientError) {
	panic(clientError)
}
