package repository

const (
	RegisterQuery        = `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)`
	CheckPasswordQuery   = `SELECT password FROM users WHERE email = $1`
	GetUserPasswordQuery = `SELECT password FROM users WHERE id = $1`
	UpdatePasswordQuery  = `UPDATE users SET password = $1 WHERE id = $2;`
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
