package repository

import (
	"auth/internal/helpers"
	"auth/internal/models"
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterUser(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req RegisterRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		hashedPassword, err := helpers.HashPassword(req.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}

		_, err = db.Exec(context.Background(), RegisterQuery, req.Name, req.Email, hashedPassword, models.UserRole)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User created successfully"})
	}
}

func LoginUser(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}
		var hashedPassword string
		err := db.QueryRow(context.Background(), CheckPasswordQuery, req.Email).Scan(&hashedPassword)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
		}

		if !helpers.CheckPassword(req.Password, hashedPassword) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}

		token, err := helpers.GenerateJWT(req.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		return c.JSON(http.StatusOK, map[string]string{"token": token})

	}
}

func ChangePassword(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Получение информации о пользователе
		user, err := helpers.GetUserByToken(c, db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		// Получение данных из запроса (старый и новый пароль)
		var req ChangePasswordRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		// Проверка на то, что старый пароль действительно принадлежит пользователю
		var hashedPassword string
		err = db.QueryRow(context.Background(), GetUserPasswordQuery, user.ID).Scan(&hashedPassword)
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid old password"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
		}

		// Сравнение старого пароля с хранимым в базе данных
		if !helpers.CheckPassword(req.OldPassword, hashedPassword) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}

		// Хэширование нового пароля
		newHashedPassword, err := helpers.HashPassword(req.NewPassword)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash new password"})
		}

		// Обновление пароля в базе данных
		_, err = db.Exec(context.Background(), UpdatePasswordQuery, newHashedPassword, user.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
		}

		// Успешное обновление
		return c.JSON(http.StatusOK, map[string]string{"message": "Password updated successfully"})
	}
}
