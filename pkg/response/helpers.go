package response

import (
	"fmt"
	"net/http"
)

func getOverrides(override []string) string {
	if len(override) > 0 {
		return override[0]
	}
	return "Resource"
}
func Created[T any](data T, resourceName ...string) GenericResponse[T] {
	name := getOverrides(resourceName)

	// Use a template for consistency
	message := fmt.Sprintf("%s created successfully", name)

	return Success(http.StatusCreated, data, message)
}

func Detail[T any](data T, resourceName ...string) GenericResponse[T] {
	// Default subject
	name := getOverrides(resourceName)
	message := fmt.Sprintf("%s retrieved successfully", name)
	return Success(http.StatusOK, data, message)
}

func List[T any](data []T, resourceName ...string) GenericResponse[[]T] {
	// Default subject
	subject := getOverrides(resourceName)
	if len(resourceName) == 0 {
		subject = "Resources"
	}
	// Use a template for consistency
	message := fmt.Sprintf("%s retrieved successfully", subject)

	return Success(http.StatusOK, data, message)
}

func Deleted(resourceName ...string) GenericResponse[struct{}] {
	// Default subject
	subject := getOverrides(resourceName)

	// Use a template for consistency
	message := fmt.Sprintf("%s deleted successfully", subject)
	return Success(http.StatusOK, struct{}{}, message)
}
