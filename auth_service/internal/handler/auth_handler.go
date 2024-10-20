package handler

import (
	"auth_service/internal/dto"
	goerrors "auth_service/internal/go_errors"
	"auth_service/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *service.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	authHandler := AuthHandler{authService: s, validator: validator.New()}
	authHandlerPtr := &authHandler
	return authHandlerPtr
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Mehtod not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Input Format", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		http.Error(w, "Request error: bad request", http.StatusBadRequest)
		log.Printf("Validation error: " + err.Error())
		return
	}

	err = h.authService.SignUp(req.Username, req.FirstName, req.LastName, req.Email, req.Role, req.Password)
	if err != nil {
		if errors.Is(err, goerrors.ErrUserAlreadyExists) {
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := make(map[string]string)
	response["Message"] = "Sign up successfull"
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var req dto.LogInRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Input Format", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		http.Error(w, "Invalid Input Format", http.StatusBadRequest)
		return
	}

	userId, token, err := h.authService.Login(req.Email, "", req.Password)
	if err != nil {
		if errors.Is(err, goerrors.ErrUserNotFound) {
			http.Error(w, "User not found, please sign up", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to login", http.StatusInternalServerError)
		}
	}

	res := dto.LogInResponse{
		UserId: userId,
		Token:  token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)

}
