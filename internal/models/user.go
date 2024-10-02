package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     uint8  `json:"role"`
	Password string `json:"password"`
}

const (
	AdminRole     = 1 // Может изменять права у всех роей кроме Admin
	ModeratorRole = 2 // Может изменять права у User
	UserRole      = 3 // Не может изменять права
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
