package response

import (
	"net/http"
)

func Created[T any](data T, messages ...string) GenericResponse[T] {
	message := "Resource created successfully"

	if len(messages) > 0 {
		message = messages[0]
	}

	return Success(http.StatusCreated, data, message)
}

func Detail[T any](data T, messages ...string) GenericResponse[T] {
	message := "Resource retrieved successfully"

	if len(messages) > 0 {
		message = messages[0]
	}

	return Success(http.StatusOK, data, message)
}

func List[T any](data []T, messages ...string) GenericResponse[[]T] {
	message := "Resources retrieved successfully"

	if len(messages) > 0 {
		message = messages[0]
	}

	return Success(http.StatusOK, data, message)
}

func Deleted() GenericResponse[struct{}] {
	return Success(http.StatusOK, struct{}{}, "Resource deleted successfully")
}
