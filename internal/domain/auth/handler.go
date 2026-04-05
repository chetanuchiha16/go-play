package auth

import (
	"encoding/json"
	"net/http"

	"github.com/chetanuchiha16/go-play/internal/api"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	apierrors "github.com/chetanuchiha16/go-play/internal/errors"
)

// Handler implements auth-related ServerInterface methods.
type Handler struct {
	service AuthService
}

// NewAuthHandler creates a new auth Handler.
func NewAuthHandler(service AuthService) *Handler {
	return &Handler{service: service}
}

// LoginUser authenticates a user and returns a JWT.
// (POST /auth/login)
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req api.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteError(w, http.StatusBadRequest, "Invalid Request Body", "Could not decode JSON request body")
		return
	}

	dbUser, token, err := h.service.Login(r.Context(), string(req.Email), req.Password)
	if err != nil {
		apiErr := apierrors.MapError(err, "user")
		api.WriteError(w, apiErr.Status, apiErr.Title, apiErr.Detail)
		return
	}

	status := "success"
	msg := "Login successful"
	createdAt := dbUser.CreatedAt.Time
	api.WriteJSON(w, http.StatusOK, api.LoginSuccessResponse{
		Status:  &status,
		Message: &msg,
		Data: &api.LoginData{
			Token: &token,
			User: &api.UserData{
				Id:        &dbUser.ID,
				Name:      &dbUser.Name,
				Email:     &dbUser.Email,
				CreatedAt: &createdAt,
			},
		},
	})
}

// NewUserResponse re-exports user.NewUserResponse for convenience.
var NewUserResponse = user.NewUserResponse
