package repository

import (
	"database/sql"
	"patient_service/internal/models"
)

type PatientRepository struct {
	DB *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	userRepo := PatientRepository{DB: db}
	userRepoPtr := &userRepo
	return userRepoPtr
}

func (r *PatientRepository) CreatePatient(patient *models.Patient) error {
	_, err := r.DB.Exec(
		"INSERT INTO patients (patient_id, first_name, last_name, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		patient.PatientID, patient.FirstName, patient.LastName, patient.Email, patient.CreatedAt, patient.UpdatedAt,
	)
	return err
}

func (r *PatientRepository) GetAllPatients() ([]*models.Patient, error) {
	rows, err := r.DB.Query(
		"SELECT patient_id, first_name, last_name, email, created_at, updated_at FROM patients",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*models.Patient
	for rows.Next() {
		patient := &models.Patient{}
		err := rows.Scan(&patient.PatientID, &patient.FirstName, &patient.LastName, &patient.Email, &patient.CreatedAt, &patient.UpdatedAt)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *PatientRepository) GetPatient(patientID string) (*models.Patient, error) {
	row := r.DB.QueryRow(
		"SELECT patient_id, first_name, last_name, email, created_at, updated_at FROM patients WHERE patient_id = $1",
		patientID,
	)

	patient := &models.Patient{}
	err := row.Scan(&patient.PatientID, &patient.FirstName, &patient.LastName, &patient.Email, &patient.CreatedAt, &patient.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return patient, nil
}

func (r *PatientRepository) RemovePatient(patientID string) error {
	_, err := r.DB.Exec(
		"DELETE FROM patients WHERE patient_id = $1",
		patientID,
	)
	return err
}

func (r *PatientRepository) UpdatePatient(patient *models.Patient) error {
	_, err := r.DB.Exec(
		"UPDATE patients SET first_name = $1, last_name = $2, email = $3, updated_at = $4 WHERE patient_id = $5",
		patient.FirstName, patient.LastName, patient.Email, patient.UpdatedAt, patient.PatientID,
	)
	return err
}

func (r *PatientRepository) CheckPatientExists(email string) (bool, error) {
	var count int
	query := `
        SELECT COUNT(*) 
    	FROM patients 
    	WHERE email = $1
    `
	err := r.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
