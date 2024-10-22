package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"patient_service/internal/dto"
	"patient_service/internal/middleware"
	"patient_service/internal/service"

	"github.com/go-playground/validator/v10"
)

type PatientHandler struct {
	patientService service.PatientServiceInterface // Use the interface here
	validator      *validator.Validate
}

func NewPatientHandler(s service.PatientServiceInterface) *PatientHandler { // Accept interface
	return &PatientHandler{patientService: s, validator: validator.New()}
}

func (h *PatientHandler) HandlePatientRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreatePatient(w, r)
	case http.MethodGet:
		h.GetPatient(w, r)
	case http.MethodPut:
		h.UpdatePatient(w, r)
	case http.MethodDelete:
		h.DeletePatient(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.ExtractClaimsFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role := claims["role"].(string)
	if role != "receptionist" {
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		return
	}

	var req dto.CreatePatientRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Input Format", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Validation error: %v", err)
		return
	}

	err = h.patientService.CreatePatient(req.FirstName, req.LastName, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error while creating patient: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "Patient created successfully"}
	json.NewEncoder(w).Encode(response)
}

func (h *PatientHandler) GetAllPatients(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.ExtractClaimsFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role := claims["role"].(string)
	if role != "doctor" && role != "receptionist" {
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		return
	}

	patients, err := h.patientService.GetAllPatients()
	if err != nil {
		http.Error(w, "Failed to fetch patients", http.StatusInternalServerError)
		log.Printf("Error while fetching patients: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(patients)
}

func (h *PatientHandler) GetPatient(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.ExtractClaimsFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role := claims["role"].(string)
	if role != "doctor" && role != "receptionist" {
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		return
	}

	patientID := r.URL.Query().Get("patient_id")
	if patientID == "" {
		http.Error(w, "Missing patient_id parameter", http.StatusBadRequest)
		return
	}

	patient, err := h.patientService.GetPatient(patientID)
	if err != nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		log.Printf("Error while fetching patient: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(patient)
}

func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.ExtractClaimsFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role := claims["role"].(string)
	if role != "doctor" && role != "receptionist" {
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		return
	}

	var req dto.UpdatePatientRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Input Format", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Validation error: %v", err)
		return
	}

	err = h.patientService.UpdatePatient(req.PatientID, req.FirstName, req.LastName, req.Email)
	if err != nil {
		http.Error(w, "Failed to update patient", http.StatusInternalServerError)
		log.Printf("Error while updating patient: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Patient updated successfully"}
	json.NewEncoder(w).Encode(response)
}

func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.ExtractClaimsFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role := claims["role"].(string)
	if role != "receptionist" {
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		return
	}

	patientID := r.URL.Query().Get("patient_id")
	if patientID == "" {
		http.Error(w, "Missing patient_id parameter", http.StatusBadRequest)
		return
	}

	err = h.patientService.RemovePatient(patientID)
	if err != nil {
		http.Error(w, "Failed to delete patient", http.StatusInternalServerError)
		log.Printf("Error while deleting patient: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Patient deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
