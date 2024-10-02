package repository

import (
	"auth/internal/helpers"
	"auth/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

const RegisterQuery = `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)`

func RegisterUser(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req RegisterRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		// Хэширование пароля
		hashedPassword, err := helpers.HashPassword(req.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}

		// Сохранение пользователя в базу данных

		_, err = db.Exec(context.Background(), RegisterQuery, req.Name, req.Email, hashedPassword, models.UserRole)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User created successfully"})
	}
}
