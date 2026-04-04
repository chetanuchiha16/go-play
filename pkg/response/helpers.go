package response

import "net/http"

func Created[T any](data T, messages ...string) {
	message := "Resource created successfully"

	if len(messages) > 0 {
		message = messages[0]
	}

	Success(http.StatusCreated, data, message)
}

func Detail[T any](data T, messages ...string) {
    message := "Resource retrieved successfully"

    if len(messages) > 0 {
        message = messages[0]
    }

    Success(http.StatusOK, data, message)
}

func List[T any](data []T, messages ...string) {
    message := "Resources retrieved successfully"

    if len(messages) > 0 {
        message = messages[0]
    }

    Success(http.StatusOK, data, message)
}

func Deleted() {
	Success(http.StatusOK, struct{}{}, "Resource deleted successfully")
}

