package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)
var (
	ErrInvalidToken          = errors.New("invalid token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
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
		return fuego.UnauthorizedError{
			Status: http.StatusBadRequest,
			Title:  "Invalid Path Parameter",
			Detail: fmt.Sprintf("The value '%s' is not a valid number", numErr.Num),
		}
	}

	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return fuego.UnauthorizedError{
			Status: http.StatusUnprocessableEntity,
			Title: "Password too long",
			Detail: "password too long",
		}
	}

	// Default to 500 for everything else
	return fuego.HTTPError{
		Status: http.StatusInternalServerError,
		Title:  "Internal Server Error",
		Detail: err.Error(),
	}
}
