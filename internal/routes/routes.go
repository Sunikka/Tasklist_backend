package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sunikka/tasklist-backendGo/internal/auth"
	"github.com/sunikka/tasklist-backendGo/internal/db"
	"github.com/sunikka/tasklist-backendGo/internal/utils"
)

type APIServer struct {
	listenAddr string
	store      db.Storage
}

func NewAPIServer(listenAddr string, store db.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s APIServer) Run() {
	mux := mux.NewRouter()

	// Request handlers
	mux.HandleFunc("/tasks/{user_id}", auth.MiddlewareJWT(createHandler(s.handleTasks), s.store))
	mux.HandleFunc("/tasks/{user_id}/{task_id}", auth.MiddlewareJWT(createHandler(s.handleTasks), s.store))

	mux.HandleFunc("/users", createHandler(s.handleUsers))
	mux.HandleFunc("/users/{user_id}", auth.MiddlewareJWT(createHandler(s.handleUsers), s.store))

	mux.HandleFunc("/login", createHandler(s.handleLogin))
	mux.HandleFunc("/register", createHandler(s.handleCreateUser))

	// TODO: CORS config
	handler := cors.Default().Handler(mux)

	log.Println("Tasklist-API listening on port", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, handler); err != nil {
		fmt.Println(err.Error())
	}
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req utils.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	if !user.ValidPassword(req.Password) {
		return fmt.Errorf("authentication failed")
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return err
	}

	response := utils.LoginResponse{
		Username: user.Name,
		Token:    token,
	}

	return utils.WriteJSON(w, 200, response)
}

// handler for  /tasks/{userID} && /tasks/{userID}/{taskID} endpoints
func (s *APIServer) handleTasks(w http.ResponseWriter, r *http.Request) error {

	// For checking if the request has :id attached to it
	_, hasID := mux.Vars(r)["task_id"]

	if !hasID {
		switch r.Method {
		case "GET":
			return s.handleGetTasksForUser(w, r)
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

// handler for  /users && /users/:id endpoints
func (s *APIServer) handleUsers(w http.ResponseWriter, r *http.Request) error {

	// For checking if the request has :id attached to it
	_, hasID := mux.Vars(r)["user_id"]

	if !hasID {
		switch r.Method {
		case "GET":
			return s.handleGetUsers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return fmt.Errorf("method not allowed")
		}

	} else {
		switch r.Method {
		case "GET":
			return s.handleGetUserByID(w, r)
		case "DELETE":
			return s.handleDeleteUser(w, r)
		case "PUT":
			return s.handleUpdateUser(w, r)
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

	return utils.WriteJSON(w, http.StatusOK, tasks)
}

func (s *APIServer) handleGetTasksForUser(w http.ResponseWriter, r *http.Request) error {
	id, err := utils.GetUserID(r)
	if err != nil {
		return err
	}

	tasks, err := s.store.GetTasksByUserID(id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, tasks)
}

func (s *APIServer) handleGetTaskByID(w http.ResponseWriter, r *http.Request) error {
	id, err := utils.GetTaskID(r)
	if err != nil {
		return err
	}

	task, err := s.store.GetTaskById(id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, task)
}

func (s *APIServer) handleCreateTask(w http.ResponseWriter, r *http.Request) error {
	req := new(utils.TaskBodyRequest)
	// user, err := auth.GetUserFromToken(r, s.store)
	userID, err := utils.GetUserID(r)
	if err != nil {
		return err
	}

	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	task, err := utils.NewTask(req.Title, req.Description, req.Deadline, userID)
	if err != nil {
		return err
	}

	created, err := s.store.CreateTask(task)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, created)
}

func (s *APIServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) error {

	id, err := utils.GetTaskID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteTask(id); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, utils.JSONres{"deleted": id})
}

func (s *APIServer) handleUpdateTask(w http.ResponseWriter, r *http.Request) error {
	req := new(utils.TaskBodyRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	id, err := utils.GetTaskID(r)
	if err != nil {
		return err
	}

	task, err := s.store.GetTaskById(id)
	if err != nil {
		return err
	}

	task.ModifyTask(req)

	if err := s.store.UpdateTask(id, task); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, utils.JSONres{"updated": id})
}

func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetUsers()
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, users)
}

func (s *APIServer) handleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	id, err := utils.GetUserID(r)
	if err != nil {
		return err
	}

	user, err := s.store.GetUserById(id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, user)

}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return fmt.Errorf("method not allowed")
	}

	req := new(utils.RegisterUserRequest)

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	user, err := utils.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateUser(user); err != nil {
		return err
	}
	return utils.WriteJSON(w, http.StatusOK, req)
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {

	id, err := utils.GetUserID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteUser(id); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, utils.JSONres{"deleted": id})
}

func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	req := new(utils.UserBodyRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	id, err := utils.GetUserID(r)
	if err != nil {
		return err
	}

	user, err := s.store.GetUserById(id)
	if err != nil {
		return err
	}

	user.ModifyUser(req)

	if err := s.store.UpdateUser(id, user); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, utils.JSONres{"updated": id})
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func createHandler(fc APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fc(w, r)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, utils.APIError{Error: err.Error()})
		}
	}
}
