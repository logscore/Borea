// All this should do is initialize a db connection, initialize handle listeners and set the listening port

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
)

var (
	URL     string
	GO_PORT string
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	defer db.DB.Close()

	http.HandleFunc("/getItems", handlers.GetItems)
	http.HandleFunc("/getItem", handlers.GetItem)
	http.HandleFunc("/createItem", handlers.CreateItem)
	http.HandleFunc("/updateItem", handlers.UpdateItem)
	http.HandleFunc("/script", handlers.HandleScriptRequest)
	http.HandleFunc("/postSession", handlers.PostSessionData)

	http.HandleFunc("/ping", handlers.PingHandler)

	GO_PORT = os.Getenv("GO_PORT")
	URL = fmt.Sprintf("0.0.0.0:%s", GO_PORT) // Changed from localhost to 0.0.0.0 for prod

	server := &http.Server{
		Addr:    URL,
		Handler: nil,
	}

	go func() {
		log.Printf("Server starting on %s\n", URL)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server init error: %v", err)
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
