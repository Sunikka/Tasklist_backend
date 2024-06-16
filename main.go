package main

// Author: Sunikka

// Anthony GG's REST API Tutorial was used as a base for migrating away from Gin
// https://www.youtube.com/watch?v=pwZuNmAzaH8&list=PL0xRBLFXXsP6nudFDqMXzrvQCZrxSOm-2

import (
	"log"
	"os"
)

func main() {
	store, err := NewStore()
	if err != nil {
		log.Fatal(err)
	}

	err = store.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	port := string(os.Getenv("SERVERPORT"))

	store.Connect()
	server := NewAPIServer(port, store)
	server.Run()

}
