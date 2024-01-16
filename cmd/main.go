package main

import (
	"fmt"
	"log"
	"test-task/db/database"
)

func main() {
	db, err := database.NewDataBase()
	if err != nil {
		log.Fatal("Error creating database connection:", err)
	}
	defer db.Close()

	fmt.Println("Connected to the database!")
}
