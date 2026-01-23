package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chetanuchiha16/go-play/db"
)

func Hello_Hina(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

type Handler struct {
	store Store
}

func NewHandler(s Store) *Handler {
	return &Handler{
		store: s,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var arg db.CreateUserParams

	error := json.NewDecoder(r.Body).Decode(&arg)

	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}

	if arg.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	user, err := h.store.CreateUser(r.Context(), arg)
	if err != nil {
		http.Error(w, "user could not be created", http.StatusInternalServerError)
		return
	}
	// CacheMutex.Lock()
	// UserCache[len(UserCache)+1] = user
	// CacheMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

	// w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.store.GetUser(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// CacheMutex.RLock()
	// user, ok := UserCache[id]
	// CacheMutex.RUnlock()

	// if !ok {
	// 	http.Error(w, "user does not exist", http.StatusNotFound)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
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

	err = h.store.DeleteUser(r.Context(), id)

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
