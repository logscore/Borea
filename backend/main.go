package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Borea/backend/db"
	"Borea/backend/handlers"

	"github.com/joho/godotenv"
)

var (
	URL    string
	goPort string
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	defer db.DB.Close()

	http.HandleFunc("/getItems", handlers.GetItems)

	err = godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	goPort = os.Getenv("GO_PORT")

	URL = fmt.Sprintf("localhost:%s", goPort)
	server := &http.Server{
		Addr:    URL,
		Handler: nil,
	}

	go func() {
		log.Printf("Server starting on %s\n", URL)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	waitForShutdown(server)
}

func waitForShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
