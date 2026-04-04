package auth

import (
	"net/http"

	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/internal/errors"
	"github.com/chetanuchiha16/go-play/pkg/response"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

type Handler struct {
	service AuthService
}

type LoginResponseWrapper response.GenericResponse[LoginResponse]

func NewAuthHandler(service AuthService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterAuthRoutes(s *fuego.Server, authmw func(http.Handler) http.Handler) {
	authRoutes := fuego.Group(s, "/auth")
	fuego.Post(authRoutes, "/login", h.Login, option.RequestContentType("application/x-www-form-urlencoded"))
}

func (h *Handler) Login(c fuego.ContextWithBody[LoginRequest]) (LoginResponseWrapper, error) {
	body, err := c.Body()
	if err != nil {
		return LoginResponseWrapper(response.Success(http.StatusBadRequest, LoginResponse{}, "Invalid request")), err
	}
	_user, token, err := h.service.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		return LoginResponseWrapper(response.Success(http.StatusUnauthorized, LoginResponse{}, "Authentication failed")), errors.MapError(err, "user")
	}

	return LoginResponseWrapper(response.Success(http.StatusOK, LoginResponse{
		Token: token,
		User:  user.NewUserResponse(_user),
	}, "Login successful")), nil
}
