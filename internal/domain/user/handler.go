package user

import (
	"net/http"
	"strconv" // To convert the ID string to an int64

	"github.com/chetanuchiha16/go-play/db"
	"github.com/chetanuchiha16/go-play/internal/errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

type Handler struct {
	service UserService
}

func NewUserHandler(s UserService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterUserRoutes(s *fuego.Server, authmw func(http.Handler) http.Handler) {

	userRoutes := fuego.Group(s, "/users")
	fuego.Post(userRoutes, "/", h.CreateUser, option.RequestContentType("application/x-www-form-urlencoded"))
	fuego.Get(userRoutes, "/", h.ListUser, option.QueryInt("limit", "Maximum number of users to return", param.Default(20)))
	fuego.Get(userRoutes, "/{id}", h.GetUser)

	//protected routes
	authGroup := fuego.Group(s, "/users", 
        option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}),
    )
	fuego.Use(authGroup, authmw)

	// Use option.Security to tell Swagger this specific route needs the token
	fuego.Delete(authGroup, "/{id}", h.DeleteUser)

}

// 1. CreateUser (STAYS THE SAME - This one uses a Body)
func (h *Handler) CreateUser(c fuego.ContextWithBody[CreateUserShema]) (UserResponse, error) {
	body, err := c.Body()
	if err != nil {
		return UserResponse{}, err
	}
	user, err := h.service.CreateUser(c.Context(), body)
	if err != nil {
		return UserResponse{}, errors.MapError(err, body.Name)
	}
	return NewUserResponse(user), errors.MapError(err, user.Name)
}

// 2. GetUser (UPDATED: Use ContextNoBody)
func (h *Handler) GetUser(c fuego.ContextNoBody) (UserResponse, error) {
	// Fuego gives you path parameters as strings
	idStr := c.PathParam("id")

	// Convert string "123" to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return UserResponse{}, errors.MapError(err, idStr) // Fuego will turn this into a 400 Bad Request automatically
	}

	user, err := h.service.GetUser(c.Context(), id)
	if err != nil {
		return UserResponse{}, errors.MapError(err, idStr)
	}

	return NewUserResponse(user), nil
}

// 3. ListUser (STAYS THE SAME)
func (h *Handler) ListUser(c fuego.ContextNoBody) ([]db.ListUsersRow, error) {
	limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 32)
	if err != nil {
		return []db.ListUsersRow{}, errors.MapError(err, "limit")
	}
	users, err := h.service.ListUsers(c.Context(), int32(limit))
	if err != nil {
		return []db.ListUsersRow{}, errors.MapError(err, "users")
	}
	return users, nil
}

// 4. DeleteUser (UPDATED: Use ContextNoBody)
func (h *Handler) DeleteUser(c fuego.ContextNoBody) (any, error) {
	idStr := c.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return db.User{}, errors.MapError(err, idStr)
	}

	err = h.service.DeleteUser(c.Context(), id)
	if err != nil {
		return nil, errors.MapError(err, idStr)
	}

	return map[string]string{"message": "user deleted"}, nil
}
