package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"os"
	"test-task/config"
	"test-task/db/database"
	"test-task/person/handler"
	"test-task/person/repository"
	"test-task/person/service"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("error to load configs: %s", err)
	}

	r := gin.Default()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := database.NewDataBase()
	if err != nil {
		log.Fatal("Error creating database connection:", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db.GetDB())
	serv := service.NewService(repo, logger.Sugar())
	hand := handler.NewHandler(serv)

	r.GET("/persons", hand.GetPeople)
	r.GET("/persons/:id", hand.GetPersonById)
	r.POST("/persons", hand.AddPerson)
	r.PUT("/persons/:id", hand.UpdatePerson)
	r.DELETE("/persons/:id", hand.DeletePerson)

	port := os.Getenv("PORT")
	err = r.Run(":" + port)

	if err != nil {
		log.Fatal("Error starting the server:", err)
	}

	fmt.Println("Connected to the database!")
}
