package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type APIError struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ResponsePermDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, APIError{Error: "permission denied"})
}

func GetUserID(r *http.Request) (uuid.UUID, error) {
	idStr := mux.Vars(r)["user_id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid user ID: %s", idStr)
	}

	return id, nil
}

func GetTaskID(r *http.Request) (uuid.UUID, error) {
	idStr := mux.Vars(r)["task_id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid task ID: %s", idStr)
	}

	return id, nil
}
