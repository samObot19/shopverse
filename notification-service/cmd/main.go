package main

import (
    "log"
    "net/http"
    "github.com/samObot19/shopverse/notification-service/internal/handlers"
    "github.com/samObot19/shopverse/notification-service/pkg/logger"
)

func main() {
    logger.Init()

    // Set up routes
    http.HandleFunc("/notify", handlers.HandleNotification)

    // Start the HTTP server
    log.Println("Starting notification service on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}