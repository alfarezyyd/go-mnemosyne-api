package exception

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
)

func ParseGormError(err error) *ClientError {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return &ClientError{
			Message:    "Record not found",
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return &ClientError{
			Message:    "Data already exists",
			StatusCode: http.StatusConflict,
		}

	// Handle MySQL/Postgres specific errors
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return &ClientError{
			Message:    "Related record not found",
			StatusCode: http.StatusBadRequest,
		}

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return &ClientError{
			Message:    "Duplicate entry",
			StatusCode: http.StatusConflict,
		}
	case errors.Is(err, gorm.ErrInvalidData):
		return &ClientError{
			Message:    "Invalid data",
			StatusCode: http.StatusBadRequest,
		}
	default:
		return &ClientError{
			Message:    "Database error occurred",
			StatusCode: http.StatusInternalServerError,
		}
	}
}
