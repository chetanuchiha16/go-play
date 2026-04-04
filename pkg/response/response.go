package response

import "time"

type GenericResponse[T any] struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    *T          `json:"data,omitempty"`
	Errors  any `json:"errors,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	RequestID string    `json:"requestId,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Page      int       `json:"page,omitempty"`
	PerPage   int       `json:"per_page,omitempty"`
	Total     int       `json:"total,omitempty"`
	Next      string    `json:"next,omitempty"`
	Prev      string    `json:"prev,omitempty"`
}

// Success helper for common success cases
func Success[T any](code int, data T, message string) GenericResponse[T] {
	resp := GenericResponse[T]{
		Status:  "success",
		Message: message,
		Data:    &data,
	}
	return resp
}

// SuccessWithMeta for responses with pagination etc.
func SuccessWithMeta[T any](code int, data T, meta Meta) GenericResponse[T]{
	resp := GenericResponse[T]{
		Status: "success",
		Data:   &data,
		Meta:   &meta,
	}
	return resp
}