// internal/erros/maperrors.go
package errors

import "github.com/go-fuego/fuego"

func MapError(err error) error {
	if err == nil {
		return nil
	}
	// Logic: If the DB says 'no rows', we tell the user 'Not Found'
	if err.Error() == "no rows in result set" {
		return fuego.NotFoundError{Title: "Resource not found"}
	}
	return fuego.InternalServerError{Title: "Internal Server Error", Detail: err.Error()}
}
