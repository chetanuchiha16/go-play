package server

import (
	"net/http"

	"github.com/chetanuchiha16/go-play/internal/api"
	"github.com/chetanuchiha16/go-play/internal/domain/auth"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Ensure Server implements the generated ServerInterface at compile time.
var _ api.ServerInterface = (*Server)(nil)

// Server composes domain handlers to satisfy api.ServerInterface.
//
// user.Handler is embedded (promotes CreateUser, GetUser, ListUsers, DeleteUser).
// auth.Handler's LoginUser is forwarded explicitly because Go disallows
// two anonymous fields with the same type name ("Handler").
type Server struct {
	*user.Handler
	authHandler *auth.Handler
	pool        *pgxpool.Pool
}

// NewServer constructs a Server with all required dependencies.
func NewServer(
	userHandler *user.Handler,
	authHandler *auth.Handler,
	pool *pgxpool.Pool,
) *Server {
	return &Server{
		Handler:     userHandler,
		authHandler: authHandler,
		pool:        pool,
	}
}

// LoginUser forwards to the auth domain handler.
func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	s.authHandler.LoginUser(w, r)
}

// HealthCheck checks database connectivity.
// (GET /health) — not domain-specific, lives at the server level.
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.pool.Ping(r.Context()); err != nil {
		api.WriteError(w, http.StatusServiceUnavailable, "Database connection failed", "Database unavailable at the moment")
		return
	}
	status := "success"
	msg := "OK"
	api.WriteJSON(w, http.StatusOK, api.HealthResponse{Status: &status, Message: &msg})
}
