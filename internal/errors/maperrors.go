package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrAlreadyExists           = errors.New("Already Exists")
	ErrNotFound                = errors.New("Not Found")
)

// APIError carries the information needed to write an RFC 7807 error response.
type APIError struct {
	Status int
	Title  string
	Detail string
}

// Error implements the error interface so APIError can be returned as an error.
func (e APIError) Error() string {
	return fmt.Sprintf("%d %s: %s", e.Status, e.Title, e.Detail)
}

// MapError translates internal errors into APIError values.
// 'res' is the name of the thing that wasn't found (e.g., "User").
func MapError(err error, res string) APIError {
	if err == nil {
		return APIError{Status: http.StatusOK}
	}

	// Database "not found" error
	if errors.Is(err, pgx.ErrNoRows) {
		return APIError{
			Status: http.StatusNotFound,
			Title:  fmt.Sprintf("%s not found", res),
			Detail: fmt.Sprintf("The requested %s does not exist in our records.", res),
		}
	}

	// Wrong password
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return APIError{
			Status: http.StatusUnauthorized,
			Title:  "Authentication Failed",
			Detail: "Incorrect Password",
		}
	}

	// Invalid numeric parameter
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		return APIError{
			Status: http.StatusBadRequest,
			Title:  "Invalid Path Parameter",
			Detail: fmt.Sprintf("The value '%s' is not a valid number", numErr.Num),
		}
	}

	// Password too long for bcrypt
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return APIError{
			Status: http.StatusUnprocessableEntity,
			Title:  "Password too long",
			Detail: "password too long",
		}
	}

	// Postgres errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return APIError{
				Status: http.StatusConflict,
				Title:  "Conflict",
				Detail: fmt.Sprintf("A %s with this unique identifier already exists.", res),
			}
		}
	}

	// Default to 500
	return APIError{
		Status: http.StatusInternalServerError,
		Title:  "Internal Server Error",
		Detail: err.Error(),
	}
}
