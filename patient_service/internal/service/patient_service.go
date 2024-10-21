package service

import (
	"log"
	goerrors "patient_service/internal/go_errors"
	"patient_service/internal/models"
	"patient_service/internal/repository"
	"patient_service/internal/utils"
)

type PatientService struct {
	repo *repository.PatientRepository
}

func NewPatientService(r *repository.PatientRepository) *PatientService {
	PatientService := PatientService{repo: r}
	PatientServicePrt := &PatientService
	return PatientServicePrt
}

func (s *PatientService) CreatePatient(firsName, lastName, email string) error {
	PatientId := utils.GetUuid()

	currentTime := utils.GetCurrentUnixTime()
	patient := models.Patient{
		PatientID: PatientId,
		FirstName: firsName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: int(currentTime),
		UpdatedAt: int(currentTime),
	}

	patientPtr := &patient

	exists, err := s.repo.CheckPatientExists(patient.Email)
	if err != nil {
		return err
	}
	if exists {
		return goerrors.ErrPatientAlreadyExists
	}

	err = s.repo.CreatePatient(patientPtr)
	if err != nil {
		log.Printf("Error occured while creating user: %v", err)
		return err
	}
	return nil
}

func (s *PatientService) GetAllPatients() ([]*models.Patient, error) {
	patients, err := s.repo.GetAllPatients()
	if err != nil {
		log.Printf("Error occurred while retrieving patients: %v", err)
		return nil, goerrors.ErrDatabase
	}
	return patients, nil
}

func (s *PatientService) GetPatient(patientID string) (*models.Patient, error) {

	patient, err := s.repo.GetPatient(patientID)
	if err != nil {
		log.Printf("Error occurred while retrieving patient: %v", err)
		return nil, goerrors.ErrDatabase
	}
	if patient == nil {
		return nil, goerrors.ErrPatientNotFound
	}
	return patient, nil
}

func (s *PatientService) UpdatePatient(patientID, firstName, lastName, email string) error {

	currentTime := utils.GetCurrentUnixTime()
	patient := models.Patient{
		PatientID: patientID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		UpdatedAt: int(currentTime),
	}

	existingPatient, err := s.repo.GetPatient(patientID)
	if err != nil {
		return goerrors.ErrDatabase
	}
	if existingPatient == nil {
		return goerrors.ErrPatientNotFound
	}

	err = s.repo.UpdatePatient(&patient)
	if err != nil {
		log.Printf("Error occurred while updating patient: %v", err)
		return goerrors.ErrDatabase
	}
	return nil
}

func (s *PatientService) RemovePatient(patientID string) error {
	patient, err := s.repo.GetPatient(patientID)
	if err != nil {
		return goerrors.ErrDatabase
	}
	if patient == nil {
		return goerrors.ErrPatientNotFound
	}

	err = s.repo.RemovePatient(patientID)
	if err != nil {
		log.Printf("Error occurred while deleting patient: %v", err)
		return goerrors.ErrDatabase
	}
	return nil
}
