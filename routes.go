package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s APIServer) Run() {
	mux := mux.NewRouter()

	// Request handlers
	mux.HandleFunc("/tasks", createHandler(s.handleTasks))
	mux.HandleFunc("/tasks/{id}", createHandler(s.handleTasks))

	// TODO: CORS config
	handler := cors.Default().Handler(mux)

	log.Println("Tasklist-API listening on port", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, handler); err != nil {
		fmt.Println(err.Error())
	}
}

// handler for  /tasks && /tasks/:id endpoints
func (s *APIServer) handleTasks(w http.ResponseWriter, r *http.Request) error {

	// For checking if the request has :id attached to it
	_, hasID := mux.Vars(r)["id"]

	if !hasID {
		switch r.Method {
		case "GET":
			return s.handleGetTasks(w, r)
		case "POST":
			return s.handleCreateTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return fmt.Errorf("method not allowed")
		}

	} else {
		switch r.Method {
		case "GET":
			return s.handleGetTaskByID(w, r)
		case "DELETE":
			return s.handleDeleteTask(w, r)
		case "PUT":
			return s.handleUpdateTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return fmt.Errorf("method not allowed")
		}
	}
}

func (s *APIServer) handleGetTasks(w http.ResponseWriter, r *http.Request) error {
	tasks, err := s.store.GetTasks()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, tasks)
}

func (s *APIServer) handleGetTaskByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	task, err := s.store.GetTaskById(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, task)
}

func (s *APIServer) handleCreateTask(w http.ResponseWriter, r *http.Request) error {
	req := new(TaskBodyRequest)

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	task := NewTask(req.Title, req.Description, req.Deadline)

	if err := s.store.CreateTask(task); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, req)
}

func (s *APIServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteTask(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, jsonRes{"deleted": id})
}

// TODO: Figure out a smarter way to handle partial objects
// Currently fetches the full object from the database, and appends the data from request into it
func (s *APIServer) handleUpdateTask(w http.ResponseWriter, r *http.Request) error {
	req := new(TaskBodyRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	id, err := getID(r)
	if err != nil {
		return err
	}

	task, err := s.store.GetTaskById(id)
	if err != nil {
		return err
	}

	task.ModifyTask(req)

	if err := s.store.UpdateTask(r, task); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, jsonRes{"updated": id})
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func createHandler(fc APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fc(w, r)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid ID: %s", idStr)
	}

	return id, nil
}
