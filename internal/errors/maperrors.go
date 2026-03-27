package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-fuego/fuego"
	"github.com/jackc/pgerrcode" // Correct
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

// MapError translates internal errors into clean API responses.
// 'res' is the name of the thing that wasn't found (e.g., "User").
func MapError(err error, res string) error {
	if err == nil {
		return nil
	}

	// Check for the specific database "not found" error
	// if err.Error() == pgx.ErrNoRows.Error() {
	if errors.Is(err, pgx.ErrNoRows) {
		return fuego.HTTPError{
			Status: http.StatusNotFound,
			Title:  fmt.Sprintf("%s not found", res),
			Detail: fmt.Sprintf("The requested %s does not exist in our records.", res),
		}
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return fuego.UnauthorizedError{
			Status: http.StatusUnauthorized,
			Title:  "Authentication Failed",
			Detail: "Incorrect Password",
		}
	}

	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		return fuego.HTTPError{
			Status: http.StatusBadRequest,
			Title:  "Invalid Path Parameter",
			Detail: fmt.Sprintf("The value '%s' is not a valid number", numErr.Num),
		}
	}

	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return fuego.UnauthorizedError{
			Status: http.StatusUnprocessableEntity,
			Title:  "Password too long",
			Detail: "password too long",
		}
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return fuego.HTTPError{
				Status: http.StatusConflict, // 409 is standard for "Already Exists"
				Title:  "Conflict",
				Detail: fmt.Sprintf("A %s with this unique identifier already exists.", res),
			}

		}
	}

	// Default to 500 for everything else
	return fuego.HTTPError{
		Status: http.StatusInternalServerError,
		Title:  "Internal Server Error",
		Detail: err.Error(),
	}
}
