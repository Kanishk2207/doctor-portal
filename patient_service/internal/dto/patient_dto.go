package dto

// CreatePatientRequest defines the structure for creating a patient
type CreatePatientRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

// UpdatePatientRequest defines the structure for updating patient information
type UpdatePatientRequest struct {
	PatientID string `json:"patient_id" validate:"required,uuid"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type PatientResponse struct {
	PatientID string `json:"patient_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}
