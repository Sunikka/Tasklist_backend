package main

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"min=5, max=30"`
	Description string `json:"description" validate:"max=100"`
	Deadline    string `json:"deadline"`
}

// Contains task fields without the ID
type TaskBodyRequest struct {
	Title       string `json:"title" validate:"min=5, max=30"`
	Description string `json:"description" validate:"max=100"`
	Deadline    string `json:"deadline"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTask(title, description, deadline string) *Task {
	return &Task{
		Title:       title,
		Description: description,
		Deadline:    deadline,
	}
}

func (t *Task) ModifyTask(req *TaskBodyRequest) {
	if req.Title != "" {
		t.Title = req.Title
	}

	if req.Description != "" {
		t.Description = req.Description
	}

	if req.Deadline != "" {
		t.Deadline = req.Deadline
	}
}

type jsonRes map[string]int
