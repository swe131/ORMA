// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"orma/db"
	"orma/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Connect to MongoDB
	if err := db.ConnectDB(); err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer db.DisconnectDB()

	// Set up CORS middleware
	corsHandler := cors.Default().Handler

	// Create a new router and register handlers
	router := mux.NewRouter()
	router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.Handle("/add-subscription", corsHandler(http.HandlerFunc(handlers.AddSubscriptionHandler))).Methods("POST")
	router.Handle("/get-subscriptions", corsHandler(http.HandlerFunc(handlers.GetSubscriptionsHandler))).Methods("GET")
	router.Handle("/add-spend", corsHandler(http.HandlerFunc(handlers.AddSpendHandler))).Methods("POST")
	router.Handle("/get-graph-data", corsHandler(http.HandlerFunc(handlers.GetGraphDataHandler))).Methods("GET")

	// Start HTTP server with CORS middleware and router
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		fmt.Println("Server listening on :8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error starting server:", err)
		}
	}()

	// Graceful shutdown on interrupt or terminate signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down server...")

	if err := server.Shutdown(nil); err != nil {
		log.Fatal("Error shutting down server:", err)
	}

	fmt.Println("Server gracefully stopped.")
}
