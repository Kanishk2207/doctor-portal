package service

import (
	goerrors "auth_service/internal/go_errors"
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"auth_service/internal/utils"
	"log"
)

type AuthService struct {
	repo *repository.UserRepositoy
}

func NewAuthService(r *repository.UserRepositoy) *AuthService {
	authService := AuthService{repo: r}
	authServicePrt := &authService
	return authServicePrt
}

func (s *AuthService) SignUp(username, firsName, lastName, email, role, password string) error {
	userId := utils.GetUuid()
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Error occured while hashing: %v", err)
	}
	currentTime := utils.GetCurrentUnixTime()
	user := models.User{
		UserID:    userId,
		Username:  username,
		FirstName: firsName,
		LastName:  lastName,
		Email:     email,
		Role:      role,
		Password:  passwordHash,
		CreatedAt: int(currentTime),
		UpdatedAt: int(currentTime),
	}

	userPtr := &user

	exists, err := s.repo.CheckUserExists(user.Email, "")
	if err != nil {
		return err
	}
	if exists {
		return goerrors.ErrUserAlreadyExists
	}

	err = s.repo.CreateUser(userPtr)
	if err != nil {
		log.Printf("Error occured while creating user: %v", err)
		return err
	}
	return nil
}

func (s *AuthService) Login(email, username, password string) (userId string, token string, err error) {
	exists, err := s.repo.CheckUserExists(email, username)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return "", "", err
	}

	if !exists {
		log.Printf("User does not exists")
		return "", "", goerrors.ErrUserNotFound
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return "", "", err
	}

	err = utils.VeryfyPassword(user.Password, password)
	if err != nil {
		log.Printf("Error occured, Password incorrect: %v", err)
	}

	token, err = utils.GenerateJwtToken(user)
	if err != nil {
		log.Fatalf("Error occured, error in token generation: %v", err)
	}
	return user.UserID, token, nil
}
