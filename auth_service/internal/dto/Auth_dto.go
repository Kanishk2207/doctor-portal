package dto

type SignUpRequest struct {
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Role      string `json:"role" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type LogInRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" valiadate:"required"`
}

type LogInResponse struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}
