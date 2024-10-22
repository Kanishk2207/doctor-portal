package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"patient_service/internal/dto"
	"patient_service/internal/handler"
	"patient_service/internal/middleware"
	"patient_service/internal/models"
	"patient_service/internal/service"
	"patient_service/internal/utils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePatient_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPatientService := service.NewMockPatientServiceInterface(ctrl)
	h := handler.NewPatientHandler(mockPatientService)

	// Mock JWT claims (role: receptionist)
	claims := map[string]interface{}{
		"role": "receptionist",
	}
	reqBody := dto.CreatePatientRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	// Create request
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8082/health", bytes.NewBuffer(reqJSON))
	req = middleware.AttachClaimsToContext(req, claims)

	w := httptest.NewRecorder()

	mockPatientService.EXPECT().CreatePatient(reqBody.FirstName, reqBody.LastName, reqBody.Email).Return(nil)

	// Execute handler
	h.CreatePatient(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Patient created successfully", resp["message"])
}

func TestCreatePatient_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPatientService := service.NewMockPatientServiceInterface(ctrl)
	h := handler.NewPatientHandler(mockPatientService)

	// Missing claims
	reqBody := dto.CreatePatientRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/patient", bytes.NewBuffer(reqJSON))
	w := httptest.NewRecorder()

	// Execute handler
	h.CreatePatient(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetPatient_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPatientService := service.NewMockPatientServiceInterface(ctrl)
	h := handler.NewPatientHandler(mockPatientService)

	claims := map[string]interface{}{
		"role": "doctor",
	}
	patientID := utils.GetUuid()
	patient := &models.Patient{
		PatientID: patientID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		CreatedAt: 1729585729,
		UpdatedAt: 1729585729,
	}

	req := httptest.NewRequest(http.MethodGet, "/patient?patient_id="+patientID, nil)
	req = middleware.AttachClaimsToContext(req, claims)
	w := httptest.NewRecorder()

	mockPatientService.EXPECT().GetPatient(patientID).Return(patient, nil)

	h.GetPatient(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Patient
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, patientID, resp.PatientID)
}

func TestGetPatient_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPatientService := service.NewMockPatientServiceInterface(ctrl)
	h := handler.NewPatientHandler(mockPatientService)

	claims := map[string]interface{}{
		"role": "doctor",
	}
	patientID := utils.GetUuid()

	req := httptest.NewRequest(http.MethodGet, "/patient?patient_id="+patientID, nil)
	req = middleware.AttachClaimsToContext(req, claims)
	w := httptest.NewRecorder()

	mockPatientService.EXPECT().GetPatient(patientID).Return(nil, errors.New("patient not found"))

	h.GetPatient(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdatePatient_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPatientService := service.NewMockPatientServiceInterface(ctrl)
	h := handler.NewPatientHandler(mockPatientService)

	claims := map[string]interface{}{
		"role": "doctor",
	}
	reqBody := dto.UpdatePatientRequest{
		PatientID: utils.GetUuid(),
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
	}

	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/patient", bytes.NewBuffer(reqJSON))
	req = middleware.AttachClaimsToContext(req, claims)
	w := httptest.NewRecorder()

	mockPatientService.EXPECT().UpdatePatient(reqBody.PatientID, reqBody.FirstName, reqBody.LastName, reqBody.Email).Return(nil)

	h.UpdatePatient(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Patient updated successfully", resp["message"])
}

func TestDeletePatient_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPatientService := service.NewMockPatientServiceInterface(ctrl)
	h := handler.NewPatientHandler(mockPatientService)

	claims := map[string]interface{}{
		"role": "receptionist",
	}
	patientID := utils.GetUuid()

	req := httptest.NewRequest(http.MethodDelete, "/patient?patient_id="+patientID, nil)
	req = middleware.AttachClaimsToContext(req, claims)
	w := httptest.NewRecorder()

	mockPatientService.EXPECT().RemovePatient(patientID).Return(nil)

	h.DeletePatient(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Patient deleted successfully", resp["message"])
}

func TestGetAllPatients_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPatientService := service.NewMockPatientServiceInterface(ctrl)
	h := handler.NewPatientHandler(mockPatientService)

	claims := map[string]interface{}{
		"role": "receptionist",
	}
	patients := []*models.Patient{
		{PatientID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
		{PatientID: "2", FirstName: "Jane", LastName: "Doe", Email: "jane.doe@example.com"},
	}

	req := httptest.NewRequest(http.MethodGet, "/patient/all", nil)
	req = middleware.AttachClaimsToContext(req, claims)
	w := httptest.NewRecorder()

	mockPatientService.EXPECT().GetAllPatients().Return(patients, nil)

	h.GetAllPatients(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []models.Patient
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "John", resp[0].FirstName)
}
