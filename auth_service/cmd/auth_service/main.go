package main

import (
	"auth_service/internal/configs"
	"auth_service/internal/db"
	"auth_service/internal/handler"
	"auth_service/internal/repository"
	"auth_service/internal/service"
	"fmt"
	"log"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "up"}`)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	config := configs.LoadConfig()
	// dsn := "dev_user:Kanishk_22@tcp(mysql-service:3306)/auth_service_db"
	// db.InitDB(dsn)
	println(config.DB_DSN)
	db.InitDB(config.DB_DSN)
	defer db.DB.Close()

	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("../../swagger-ui/"))))

	http.HandleFunc("/health", healthCheckHandler)
	port := config.HTTPAddress

	userRepo := repository.NewUserRepository(db.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	http.HandleFunc("/signup", authHandler.Signup)
	http.HandleFunc("/login", authHandler.Login)

	fmt.Printf("Starting server on port %s...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("Could not start server: %v", err)
	}
}
