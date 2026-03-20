package errors

import (
	"fmt"
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5"
)

// MapError translates internal errors into clean API responses.
// 'res' is the name of the thing that wasn't found (e.g., "User").
func MapError(err error, res string) error {
	if err == nil {
		return nil
	}

	// Check for the specific database "not found" error
	if err.Error() == pgx.ErrNoRows.Error() {
		return fuego.HTTPError{
			Status: http.StatusNotFound,
			Title:  fmt.Sprintf("%s not found", res),
			Detail: fmt.Sprintf("The requested %s does not exist in our records.", res),
		}
	}

	// Default to 500 for everything else
	return fuego.HTTPError{
		Status: http.StatusInternalServerError,
		Title:  "Internal Server Error",
		Detail: err.Error(),
	}
}
