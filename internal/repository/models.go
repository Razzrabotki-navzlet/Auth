package repository

const (
	RegisterQuery       = `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)`
	CheckPasswordQuery  = `SELECT password FROM users WHERE email = $1`
	UpdatePasswordQuery = `UPDATE users SET password = $1 WHERE email = $2;`
	UpdateRoleQuery     = `UPDATE users SET role = $1 WHERE id = $2`
	CheckVerify         = `SELECT is_verified FROM users WHERE email = $1`
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

type ResetPasswordRequest struct {
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
	Token           string `json:"token"`
}

type ChangeRoleRequest struct {
	UserId      int `json:"user_id"`
	CurrentRole int `json:"current_role"`
	NewRole     int `json:"new_role"`
}
