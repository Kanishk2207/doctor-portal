package main

import (
	"fmt"
	"log"
	"net/http"
	"patient_service/internal/configs"
	"patient_service/internal/db"
	"patient_service/internal/handler"
	"patient_service/internal/middleware"
	"patient_service/internal/repository"
	"patient_service/internal/service"
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

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthCheckHandler)

	patientRepo := repository.NewPatientRepository(db.DB)
	patientService := service.NewPatientService(patientRepo)
	patientHandler := handler.NewPatientHandler(patientService)

	mux.HandleFunc("/patient", patientHandler.HandlePatientRoutes)
	mux.HandleFunc("/patient/all", patientHandler.GetAllPatients)

	protectedHandler := middleware.TokenValidationMiddleware(mux)

	port := config.HTTPAddress
	fmt.Printf("Starting server on port %s...\n", port)
	err := http.ListenAndServe(port, protectedHandler)
	if err != nil {
		log.Printf("Could not start server: %v", err)
	}
}
