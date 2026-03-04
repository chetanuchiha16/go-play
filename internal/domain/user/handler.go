package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chetanuchiha16/go-play/db"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func Hello_Hina(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

type Handler struct {
	svc      Service
	validate *validator.Validate
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc:      svc,
		validate: validator.New(),
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserShema

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	arg := db.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	// if arg.Name == "" {
	// 	http.Error(w, "Name is required", http.StatusBadRequest)
	// 	return
	// }

	user, err := h.svc.CreateUser(r.Context(), arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	// CacheMutex.Lock()
	// UserCache[len(UserCache)+1] = user
	// CacheMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

	// w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.svc.GetUser(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	res := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	// CacheMutex.RLock()
	// user, ok := UserCache[id]
	// CacheMutex.RUnlock()

	// if !ok {
	// 	http.Error(w, "user does not exist", http.StatusNotFound)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	// data, err := json.Marshal(user)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	// w.Write(data)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	// claims, err := ValidateToken(token)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	// if claims["user_id"] != id {
	// 	http.Error(w, "You cannot delete other user", http.StatusUnauthorized)
	// 	return
	// }

	auth_id := int64(r.Context().Value("user_id").(float64)) // In JSON, all numbers are floats. Since JWTs are just encoded JSON, the library converts your user_id to a float64 by default.
	if auth_id != id {
		http.Error(w, "You cannot delete other user", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	// CacheMutex.Lock()
	// defer CacheMutex.Unlock()
	// if _, ok := UserCache[id]; !ok {
	// 	http.Error(w, "user does not exist", http.StatusNotFound)
	// 	return
	// }

	err = h.svc.DeleteUser(r.Context(), id)

	if err != nil {
		http.Error(w, "could not be deleted", http.StatusInternalServerError)
		return
	}

	// CacheMutex.Lock()
	// delete(UserCache, id)

	// _, err := json.Marshal(user)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var emailAndPassword struct{ Email, Password string } // capital so json can see it
	json.NewDecoder(r.Body).Decode(&emailAndPassword)
	user, token, err := h.svc.Login(r.Context(), emailAndPassword.Email, emailAndPassword.Password)
	if err != nil {
		log.Error().Err(err).Msg("Login failed")
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	res := LoginResponse{
		Token: token,
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
