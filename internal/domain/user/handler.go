package user

import (
	"net/http"
	"strconv" // To convert the ID string to an int64

	"github.com/chetanuchiha16/go-play/db"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h Handler) RegisterUserRoutes(s *fuego.Server, authmw func(http.Handler) http.Handler) {
	fuego.Post(s, "/login", h.Login)
	fuego.Post(s, "/users", h.CreateUser)
	
	userRoutes := fuego.Group(s, "/users")
	fuego.Get(userRoutes, "/s", h.ListUser)
	fuego.Get(userRoutes, "/{id}", h.GetUser)

	authGroup := fuego.Group(s, "")
	fuego.Use(authGroup, authmw)

	// Use option.Security to tell Swagger this specific route needs the token
	fuego.Delete(authGroup, "/users/{id}", h.DeleteUser, option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}))

}

// 1. CreateUser (STAYS THE SAME - This one uses a Body)
func (h *Handler) CreateUser(c fuego.ContextWithBody[db.CreateUserParams]) (db.User, error) {
	body, err := c.Body()
	if err != nil {
		return db.User{}, err
	}
	return h.service.CreateUser(c.Context(), body)
}

// 2. GetUser (UPDATED: Use ContextNoBody)
func (h *Handler) GetUser(c fuego.ContextNoBody) (db.User, error) {
	// Fuego gives you path parameters as strings
	idStr := c.PathParam("id")
	
	// Convert string "123" to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return db.User{}, err // Fuego will turn this into a 400 Bad Request automatically
	}

	return h.service.GetUser(c.Context(), id)
}

// 3. ListUser (STAYS THE SAME)
func (h *Handler) ListUser(c fuego.ContextNoBody) ([]db.User, error) {
	return h.service.ListUsers(c.Context())
}

// 4. DeleteUser (UPDATED: Use ContextNoBody)
func (h *Handler) DeleteUser(c fuego.ContextNoBody) (any, error) {
	idStr := c.PathParam("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	err := h.service.DeleteUser(c.Context(), id)
	if err != nil {
		return nil, err
	}
	return map[string]string{"message": "user deleted"}, nil
}

// 5. Login (STAYS THE SAME - This one uses a Body)
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Use the LoginResponse struct instead of map[string]string
func (h *Handler) Login(c fuego.ContextWithBody[LoginRequest]) (LoginResponse, error) {
    body, _ := c.Body()
    user, token, err := h.service.Login(c.Context(), body.Email, body.Password)
    if err != nil {
        return LoginResponse{}, fuego.UnauthorizedError{
			Title: "you are not a user, please register",
			Detail: err.Error(),
		}
    }

    // Return a structured response that Swagger can read
    return LoginResponse{
        Token: token,
        User: UserResponse{
            ID:        user.ID,
            Name:      user.Name,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
        },
    }, nil
}