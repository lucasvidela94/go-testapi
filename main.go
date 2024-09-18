package main

import (
	"database/sql"
	"log"
	"net/http"
	"taskmanager/handler"
	"taskmanager/service"
	"taskmanager/store"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file source driver
	_ "github.com/lib/pq"                                // PostgreSQL driver
)

func main() {
	// Database connection
	dataSourceName := "postgres://taskuser:nueva_contrasena@localhost:5432/taskmanager?sslmode=disable"

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not start SQL driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Correctly specify the file scheme
		"postgres", driver)
	if err != nil {
		log.Fatalf("Could not start migration: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	} else {
		log.Println("Migrations applied successfully")
	}

	// Initialize the store
	taskStore := &store.SQLTaskStore{DB: db}
	taskService := service.TaskService{TaskStore: taskStore}
	taskHandler := handler.TaskHandler{TaskService: taskService}

	// Define routes and handlers
	http.HandleFunc("/tasks", taskHandler.CreateTaskHandler)
	http.HandleFunc("/tasks/list", taskHandler.ListTasksHandler)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
