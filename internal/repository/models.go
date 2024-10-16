package repository

const (
	RegisterQuery       = `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)`
	GetUserData         = `SELECT id, password, role FROM users WHERE email = $1`
	UpdatePasswordQuery = `UPDATE users SET password = $1 WHERE email = $2;`
	UpdateRoleQuery     = `UPDATE users SET role = $1 WHERE id = $2`
	UpdateVerifyQuery   = `UPDATE users SET is_verified=true WHERE email=$1`
	CheckVerify         = `SELECT is_verified FROM users WHERE email = $1`
	AddLogsQuery        = `INSERT INTO request_logs (user_id, method, path, status_code, created_at) VALUES ($1, $2, $3, $4, $5)`
	GetLogsQuery        = `SELECT id, user_id, method, path, status_code, created_at FROM request_logs ORDER BY created_at DESC`
	GetUserListQuery    = `SELECT id, name, email, role FROM users`
	ChangeRoleQuery     = `UPDATE users SET role = $1 WHERE id = $2`
)

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
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
