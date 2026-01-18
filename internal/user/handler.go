package user

import (
	"net/http"
	"fmt"
	"encoding/json"
	"strconv"
)

func Hello_Hina(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User

	error := json.NewDecoder(r.Body).Decode(&user)

	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	CacheMutex.Lock()
	UserCache[len(UserCache)+1] = user
	CacheMutex.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	CacheMutex.RLock()
	user, ok := UserCache[id]
	CacheMutex.RUnlock()

	if !ok {
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	CacheMutex.Lock()
	defer CacheMutex.Unlock()
	if _, ok := UserCache[id]; !ok {
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	CacheMutex.Lock()
	delete(UserCache, id)

	// _, err := json.Marshal(user)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	w.WriteHeader(http.StatusNoContent)
}
