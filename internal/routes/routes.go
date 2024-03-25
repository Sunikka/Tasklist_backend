package routes

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sunikka/tasklist-backendGo/internal/database"
)

var db *sql.DB

func NewRouter() {
	database.InitDB()

	router := gin.Default()
	// TODO: Cors config, currently allow-all-origins
	router.Use(cors.Default())

	// Request handlers
	router.GET("/tasks", database.GetTasks)
	router.GET("/tasks/:id", database.GetTaskById)
	router.POST("/tasks", database.CreateTask)
	router.DELETE("/tasks/:id", database.DeleteTask)
	router.PUT("/tasks/:id", database.UpdateTask)
	router.Run("127.0.0.1:8000")
}
