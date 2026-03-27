package auth

import (
	"net/http"

	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/internal/errors"
	"github.com/go-fuego/fuego"
)

type Handler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterAuthRoutes(s *fuego.Server, authmw func(http.Handler) http.Handler) {
	authRoutes := fuego.Group(s, "/auth")
	fuego.Post(authRoutes, "/login", h.Login)
}

// Use the LoginResponse struct instead of map[string]string
func (h *Handler) Login(c fuego.ContextWithBody[LoginRequest]) (LoginResponse, error) {
	body, err := c.Body()
	if err != nil {
		return LoginResponse{}, err
	}
	_user, token, err := h.service.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		return LoginResponse{}, errors.MapError(err, "user")
	}

	// Return a structured response that Swagger can read
	return LoginResponse{
		Token: token,
		User: user.UserResponse{
			ID:        _user.ID,
			Name:      _user.Name,
			Email:     _user.Email,
			CreatedAt: _user.CreatedAt,
		},
	}, nil
}
