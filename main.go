package main

import (
	"context"
	"fmt"
	"infilon_task/controller"
	"infilon_task/database"
	"infilon_task/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {
	//Load env
	godotenv.Load()

	// Database connection
	database.ConnectDatabase()

	//PersonService
	personService := service.NewPersonService(database.DB)

	//PersonController
	personController := controller.NewPersonController(personService)

	r := gin.Default()

	// Endpoints
	r.GET("/person/:person_id/info", personController.GetPersonInfo)
	r.POST("/person/create", personController.CreatePerson)


	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Listen termination calls
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Server started...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe error: %s\n", err)
		}
	}()

	// Wait for a termination signal
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Shutting down gracefully...")

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown error: %s\n", err)
	}

	fmt.Println("Server shutdown complete")

	// Close database connection
	database.DB.Close()
}
