package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chetanuchiha16/go-play/internal/api"
	apierrors "github.com/chetanuchiha16/go-play/internal/errors"
	"github.com/go-playground/validator/v10"
)

// Handler implements user-related ServerInterface methods.
type Handler struct {
	service  UserService
	validate *validator.Validate
}

// NewHandler creates a new user Handler.
func NewHandler(s UserService) *Handler {
	return &Handler{
		service: s, 
		validate: validator.New(),
	}
}

// CreateUser registers a new user.
// (POST /users)
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req api.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "Invalid Request Body", "Could not decode JSON request body")
		return
	}

	schema := CreateUserShema{
		Name:     req.Name,
		Email:    string(req.Email),
		Password: req.Password,
	}

	if err := h.validate.Struct(schema); err != nil {
		api.WriteError(w, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	dbUser, err := h.service.CreateUser(r.Context(), schema)
	if err != nil {
		apiErr := apierrors.MapError(err, schema.Name)
		api.WriteError(w, apiErr.Status, apiErr.Title, apiErr.Detail)
		return
	}

	status := http.StatusCreated
	msg := fmt.Sprintf("User %v created successfully", dbUser.Name)
	createdAt := dbUser.CreatedAt.Time
	api.WriteJSON(w, http.StatusCreated, api.UserSuccessResponse{
		Status:  &status,
		Message: &msg,
		Data: &api.UserData{
			Id:        &dbUser.ID,
			Name:      &dbUser.Name,
			Email:     &dbUser.Email,
			CreatedAt: &createdAt,
		},
	})
}

// ListUsers returns a paginated list of users.
// (GET /users)
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request, params api.ListUsersParams) {
	limit := int32(20)
	if params.Limit != nil {
		limit = int32(*params.Limit)
	}

	users, err := h.service.ListUsers(r.Context(), limit)
	if err != nil {
		apiErr := apierrors.MapError(err, "users")
		api.WriteError(w, apiErr.Status, apiErr.Title, apiErr.Detail)
		return
	}

	data := make([]api.UserData, 0, len(users))
	for _, u := range users {
		id := u.ID
		name := u.Name
		email := u.Email
		createdAt := u.CreatedAt.Time
		data = append(data, api.UserData{
			Id:        &id,
			Name:      &name,
			Email:     &email,
			CreatedAt: &createdAt,
		})
	}

	status := http.StatusOK
	msg := "Users retrieved successfully"
	api.WriteJSON(w, http.StatusOK, api.UserListResponse{
		Status:  &status,
		Message: &msg,
		Data:    &data,
	})
}

// GetUser retrieves a single user by ID.
// (GET /users/{id})
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, id api.UserId) {
	dbUser, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		apiErr := apierrors.MapError(err, fmt.Sprintf("%d", id))
		api.WriteError(w, apiErr.Status, apiErr.Title, apiErr.Detail)
		return
	}

	status := http.StatusOK
	msg := fmt.Sprintf("User %v retrieved successfully", dbUser.Name)
	createdAt := dbUser.CreatedAt.Time
	api.WriteJSON(w, http.StatusOK, api.UserSuccessResponse{
		Status:  &status,
		Message: &msg,
		Data: &api.UserData{
			Id:        &dbUser.ID,
			Name:      &dbUser.Name,
			Email:     &dbUser.Email,
			CreatedAt: &createdAt,
		},
	})
}

// DeleteUser removes a user by ID (requires Bearer auth).
// (DELETE /users/{id})
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, id api.UserId) {
	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		apiErr := apierrors.MapError(err, fmt.Sprintf("%d", id))
		api.WriteError(w, apiErr.Status, apiErr.Title, apiErr.Detail)
		return
	}

	status := http.StatusOK
	msg := fmt.Sprintf("User %v deleted successfully", id)
	api.WriteJSON(w, http.StatusOK, api.DeleteResponse{
		Status:  &status,
		Message: &msg,
	})
}
