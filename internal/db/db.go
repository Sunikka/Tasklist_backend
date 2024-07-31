package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sunikka/tasklist-backendGo/internal/utils"
)

type Storage interface {
	GetTasks() ([]utils.Task, error)
	GetTasksByUserID(userID uuid.UUID) ([]utils.Task, error)
	GetTaskById(id uuid.UUID) (utils.Task, error)
	CreateTask(task *utils.Task) (*utils.Task, error)
	DeleteTask(id uuid.UUID) error
	UpdateTask(id uuid.UUID, task utils.Task) error
	GetUsers() ([]utils.User, error)
	CreateUser(user *utils.User) error
	GetUserById(id uuid.UUID) (utils.User, error)
	DeleteUser(id uuid.UUID) error
	UpdateUser(id uuid.UUID, user utils.User) error
	GetUserByEmail(email string) (utils.User, error)
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

	db, err := sql.Open("mysql", connStr+"?parseTime=true")
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

func (m MySQLStore) GetTasks() ([]utils.Task, error) {
	var tasks []utils.Task

	rows, err := m.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var task utils.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Deadline,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.UserID,
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

func (m MySQLStore) GetTasksByUserID(userID uuid.UUID) ([]utils.Task, error) {
	var tasks []utils.Task
	userIDBin, err := userID.MarshalBinary()
	if err != nil {
		return nil, err
	}

	rows, err := m.db.Query("SELECT * FROM tasks WHERE user_id = ?", userIDBin)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var task utils.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Deadline,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.UserID,
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

func (m *MySQLStore) GetTaskById(id uuid.UUID) (utils.Task, error) {
	var task utils.Task
	idBin, err := id.MarshalBinary()
	if err != nil {
		return task, err
	}

	row := m.db.QueryRow("SELECT * FROM tasks WHERE task_id = ?", idBin)

	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Deadline, &task.CreatedAt, &task.UpdatedAt, &task.UserID); err != nil {
		return task, err
	}

	return task, nil
}

func (m *MySQLStore) CreateTask(task *utils.Task) (*utils.Task, error) {
	queryStr := `INSERT INTO tasks (task_id, title,description, deadline, created_at, updated_at, user_id) VALUES (?, ?, ?, ?, ?, ?, ?)`

	taskID := uuid.New()
	taskIDBin, err := taskID.MarshalBinary()
	if err != nil {
		return nil, err
	}

	userIDBin, err := task.UserID.MarshalBinary()
	if err != nil {
		return nil, err
	}

	_, err = m.db.Exec(queryStr, taskIDBin, task.Title, task.Description, task.Deadline, time.Now().UTC(), time.Now().UTC(), userIDBin)
	if err != nil {
		return nil, err
	}

	created, err := m.GetTaskById(taskID)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *MySQLStore) DeleteTask(id uuid.UUID) error {
	idBin, err := id.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = s.db.Exec("DELETE FROM tasks WHERE task_id = ?", idBin)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStore) UpdateTask(id uuid.UUID, task utils.Task) error {
	idBin, err := id.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE tasks SET title = ?, description = ?, deadline = ?, updated_at = ? WHERE task_id = ?",
		task.Title, task.Description, task.Deadline, time.Now().UTC(), idBin)
	if err != nil {
		return err
	}

	return nil

}

func (m *MySQLStore) CreateUser(user *utils.User) error {
	queryStr := `INSERT INTO users (user_id, username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?) `

	userID, err := uuid.New().MarshalBinary()
	if err != nil {
		return err
	}

	_, err = m.db.Exec(queryStr, userID, user.Name, user.Email, user.HashedPw, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}

func (m MySQLStore) GetUsers() ([]utils.User, error) {
	var users []utils.User

	rows, err := m.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user utils.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.HashedPw,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *MySQLStore) GetUserById(id uuid.UUID) (utils.User, error) {
	var user utils.User
	idBin, err := id.MarshalBinary()
	if err != nil {
		return user, err
	}

	row := m.db.QueryRow("SELECT * FROM users WHERE user_id = ?", idBin)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPw, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return user, err
	}

	return user, nil
}

func (m *MySQLStore) GetUserByEmail(email string) (utils.User, error) {
	var user utils.User

	row := m.db.QueryRow("SELECT * FROM users WHERE email = ?", email)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPw, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return user, err
	}

	return user, nil
}

func (s *MySQLStore) DeleteUser(id uuid.UUID) error {
	idBin, err := id.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM users WHERE user_id = ?", idBin)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStore) UpdateUser(id uuid.UUID, user utils.User) error {
	idBin, err := id.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE users SET username = ?, email = ?, password = ?, updated_at = ? WHERE user_id = ?",
		user.Name, user.Email, user.HashedPw, time.Now().UTC(), idBin)
	if err != nil {
		return err
	}

	return nil

}

func (s *MySQLStore) InitDB() error {
	err := s.createUsersTable()
	if err != nil {
		return err
	}
	err = s.createTasksTable()
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) createTasksTable() error {
	queryStr := `CREATE TABLE IF NOT EXISTS tasks 
	(
		task_id BINARY(16) NOT NULL ,
		title VARCHAR(255) NOT NULL,
		description VARCHAR(255) NOT NULL,
		deadline DATE NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		user_id BINARY(16) NOT NULL,
		
		PRIMARY KEY(task_id),	
		FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
		);`

	_, err := s.db.Exec(queryStr)
	return err
}

func (s *MySQLStore) createUsersTable() error {
	queryStr := `CREATE TABLE IF NOT EXISTS users (
		user_id BINARY(16) NOT NULL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);`

	_, err := s.db.Exec(queryStr)
	return err
}
