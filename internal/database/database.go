package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}

func InitDB() {
	// Load environment var (
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env")
	}
	// Database config
	DBconf := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DBSERVER"),
		DBName: os.Getenv("DBNAME"),
	}
	// Get database connection
	var err error
	db, err = sql.Open("mysql", DBconf.FormatDSN())
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

func GetTasks(c *gin.Context) {
	// c.IndentedJSON(http.StatusOK, tasks)
	var tasks []task

	rows, err := db.Query("SELECT * FROM Assignment")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var task task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Deadline); err != nil {
			panic(err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func GetTaskById(c *gin.Context) {
	// SQL version
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	var task task

	row := db.QueryRow("SELECT * FROM Assignment WHERE id = ?", id)

	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Deadline); err != nil {
		if err == sql.ErrNoRows {
			panic(err)
		}
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	var newTask task

	// Bind data to from JSON to the new task obj
	if err := c.BindJSON(&newTask); err != nil {
		return
	}

	result, err := db.Exec("INSERT INTO Assignment (title,description, deadline) VALUES (?, ?, ?)", &newTask.Title, newTask.Description, &newTask.Deadline)
	if err != nil {
		panic(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Created a new Assignment with id %v", id)})

}

func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	result, err := db.Exec("DELETE FROM Assignment WHERE id = ?", id)
	if err != nil {
		panic(err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deletion succesful, affected rows: %v", affectedRows)})

}

// TODO: PUT request handling
func UpdateTask(c *gin.Context) {
	var updatedTask task

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	row := db.QueryRow("SELECT * FROM Assignment WHERE id = ?", id)

	if err := row.Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.Deadline); err != nil {
		if err == sql.ErrNoRows {
			panic(err)
		}
		panic(err)
	}

	if err := c.BindJSON(&updatedTask); err != nil {
		panic(err)
	}

	result, err := db.Exec("UPDATE Assignment SET title = ?, description = ?, deadline = ? WHERE id = ?", &updatedTask.Title, &updatedTask.Description, &updatedTask.Deadline, id)
	if err != nil {
		panic(err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Update succesful, affected rows: %v", affectedRows)})
}
