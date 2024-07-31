package utils

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Task struct {
	ID          uuid.UUID `json:"task_id"`
	Title       string    `json:"title" validate:"min=5, max=30"`
	Description string    `json:"description" validate:"max=100"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uuid.UUID `json:"user_id"`
}

// Contains task fields without the ID
type TaskBodyRequest struct {
	Title       string    `json:"title" validate:"min=5, max=30"`
	Description string    `json:"description" validate:"max=100"`
	Deadline    string    `json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uuid.UUID `json:"user_id"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"username"`
	Email     string    `json:"email"`
	HashedPw  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserBodyRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func NewTask(title, description, deadline string, userID uuid.UUID) (*Task, error) {
	// Time of day for the deadline currently hardcoded into 23:59 PM
	dlParsed, err := time.Parse(time.RFC3339, deadline+"T23:59:00Z")
	if err != nil {
		return nil, err
	}

	return &Task{
		Title:       title,
		Description: description,
		Deadline:    dlParsed,
		UserID:      userID,
	}, nil
}
func (t *Task) ModifyTask(req *TaskBodyRequest) error {
	if req.Title != "" {
		t.Title = req.Title
	}

	if req.Description != "" {
		t.Description = req.Description
	}

	if req.Deadline != "" {
		dlParsed, err := time.Parse(time.RFC3339, req.Deadline)
		if err != nil {
			return err
		}
		t.Deadline = dlParsed
	}

	return nil
}

func (u *User) ModifyUser(req *UserBodyRequest) error {
	if req.Username != "" {
		u.Name = req.Username
	}

	if req.Email != "" {
		u.Email = req.Email
	}

	if req.Password != "" {
		hashPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		u.HashedPw = string(hashPw)
	}

	return nil
}

func NewUser(name, email, password string) (*User, error) {
	hashPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:     name,
		Email:    email,
		HashedPw: string(hashPw),
	}, nil

}

func (u *User) ValidPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPw), []byte(password))
	return err == nil
}

type JSONres map[string]uuid.UUID
