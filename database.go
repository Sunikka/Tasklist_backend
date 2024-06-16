package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Storage interface {
	GetTasks() ([]Task, error)
	GetTaskById(id int) (Task, error)
	CreateTask(task *Task) error
	DeleteTask(id int) error
	UpdateTask(r *http.Request, task Task) error
}

type MySQLStore struct {
	db *sql.DB
}

func NewStore() (*MySQLStore, error) {
	// Load environment var (
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	connStr := os.Getenv("DBUSER") + ":" + os.Getenv("DBPASS") + "@/" + os.Getenv("DBNAME")
	// fmt.Println("Connection string:", connStr)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &MySQLStore{
		db: db,
	}, nil
}

func (m *MySQLStore) Connect() {

	// Database config
	DBconf := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DBSERVER"),
		DBName: os.Getenv("DBNAME"),
	}
	// Get database connection
	db, err := sql.Open("mysql", DBconf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Attempting connection...")
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func (m MySQLStore) GetTasks() ([]Task, error) {
	var tasks []Task

	rows, err := m.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Deadline,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *MySQLStore) GetTaskById(id int) (Task, error) {
	var task Task

	row := m.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Deadline); err != nil {
		return task, err
	}

	return task, nil
}

func (m *MySQLStore) CreateTask(task *Task) error {
	queryStr := `INSERT INTO tasks (title,description, deadline) VALUES (?, ?, ?)`

	_, err := m.db.Exec(queryStr, task.Title, task.Description, task.Deadline)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStore) DeleteTask(id int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStore) UpdateTask(r *http.Request, task Task) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE tasks SET title = ?, description = ?, deadline = ? WHERE id = ?",
		task.Title, task.Description, task.Deadline, id)
	if err != nil {
		return err
	}

	return nil

}

func (s *MySQLStore) InitDB() error {
	return s.createTasksTable()
}

func (s *MySQLStore) createTasksTable() error {
	queryStr := `CREATE TABLE IF NOT EXISTS tasks 
	(
		id INT NOT NULL AUTO_INCREMENT,
		title VARCHAR(255) NOT NULL,
		description VARCHAR(255) NOT NULL,
		deadline DATE NOT NULL,
		PRIMARY KEY (id)
		);`

	_, err := s.db.Exec(queryStr)
	return err
}
