package repository

import (
	"auth_service/internal/models"
	"database/sql"
)

type UserRepositoy struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoy {
	userRepo := UserRepositoy{DB: db}
	userRepoPtr := &userRepo
	return userRepoPtr
}

func (r *UserRepositoy) CheckUserExists(email, username string) (bool, error) {
	var count int
	query := `
        SELECT COUNT(*) 
    	FROM users 
    	WHERE email = $1 OR username = $2
    `
	err := r.DB.QueryRow(query, email, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepositoy) CreateUser(user *models.User) error {
	_, err := r.DB.Exec(
		"INSERT INTO users (user_id, username, first_name, last_name, email, role, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		user.UserID, user.Username, user.FirstName, user.LastName, user.Email, user.Role, user.Password, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepositoy) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := r.DB.QueryRow(
		"SELECT user_id, username, first_name, last_name, email, role, password FROM users WHERE email = $1", email,
	).Scan(
		&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.Password,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}
