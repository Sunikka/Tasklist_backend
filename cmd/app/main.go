package main

import (
	"github.com/sunikka/tasklist-backendGo/internal/routes"
)

// TODO
// type user struct {
// 	ID       int    `json:"id"`
// 	Email    string `json:"email"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

func main() {
	routes.NewRouter()
}
