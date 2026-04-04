package response

import "net/http"

func Created[T any](data T, messages ...string) {
	message := "Resource created successfully"

	if len(messages) > 0 {
		message = messages[0]
	}

	Success(http.StatusCreated, data, message)
}

func Deleted() {
	Success(http.StatusOK, struct{}{}, "Resource deleted successfully")
}