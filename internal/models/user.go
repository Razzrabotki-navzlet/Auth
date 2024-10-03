package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     uint8  `json:"role"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  uint8  `json:"role"`
}

const (
	AdminRole   = 1
	RaterRole   = 2
	WatcherRole = 3
)
