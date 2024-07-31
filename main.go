package main

// Author: Sunikka

// Credit: Anthony GG's REST API Tutorial was used as a base for removing Gin from dependencies
// https://www.youtube.com/watch?v=pwZuNmAzaH8&list=PL0xRBLFXXsP6nudFDqMXzrvQCZrxSOm-2

import (
	"log"
	"os"

	"github.com/sunikka/tasklist-backendGo/internal/db"
	"github.com/sunikka/tasklist-backendGo/internal/routes"
)

func main() {
	store, err := db.NewStore()
	if err != nil {
		log.Fatal(err)
	}

	err = store.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	port := string(os.Getenv("SERVERPORT"))

	store.Connect()
	server := routes.NewAPIServer(port, store)
	server.Run()

}
