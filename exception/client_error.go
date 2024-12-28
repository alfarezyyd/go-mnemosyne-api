package exception

import "fmt"

type ClientError struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Trace      interface{} `json:"trace"`
}

func (clientError *ClientError) Error() string {
	return fmt.Sprintf("Error %d: %s", clientError.StatusCode, clientError.Message)
}

func NewClientError(statusCode int, message string, traceError ...interface{}) *ClientError {
	return &ClientError{
		StatusCode: statusCode,
		Message:    message,
		Trace:      traceError,
	}
}

func ThrowClientError(clientError *ClientError) {
	panic(clientError)
}
