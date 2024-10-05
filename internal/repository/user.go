package repository

import (
	"auth/internal/helpers"
	"auth/internal/models"
	"context"
	"database/sql"
	"fmt"
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

		_, err = db.Exec(context.Background(), RegisterQuery, req.Name, req.Email, hashedPassword, models.WatcherRole)
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
			fmt.Println("Error binding request:", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}
		fmt.Println("REQUEST /home/og/projects/Auth/internal/repository/user.go", req.Email)

		var hashedPassword string
		err := db.QueryRow(context.Background(), CheckPasswordQuery, req.Email).Scan(&hashedPassword)
		if err == sql.ErrNoRows {
			fmt.Println("No user found with email:", req.Email)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		} else if err != nil {
			fmt.Println("Error querying database:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
		}

		if !helpers.CheckPassword(req.Password, hashedPassword) {
			fmt.Println("Password does not match for user:", req.Email)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}

		token, err := helpers.GenerateJWT(req.Email)
		if err != nil {
			fmt.Println("Error generating JWT:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		fmt.Println("Login successful for user:", req.Email)
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

		// Получение данных из запроса (старый и новый пароль) (нужно убедиться что пользователь знает старый пароль)
		var req ChangePasswordRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		// Проверка на то, что старый пароль действительно принадлежит пользователю
		var oldPassword string
		err = db.QueryRow(context.Background(), CheckPasswordQuery, user.Email).Scan(&oldPassword)
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid old password"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
		}

		// Сравнение нового пароля с хранимым в базе данных, чтоб не совпадал
		if helpers.CheckPassword(req.NewPassword, oldPassword) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "New password should not repeat old password"})
		}

		// Хэширование нового пароля
		newHashedPassword, err := helpers.HashPassword(req.NewPassword)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash new password"})
		}

		// Обновление пароля в базе данных
		_, err = db.Exec(context.Background(), UpdatePasswordQuery, newHashedPassword, user.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
		}

		// Успешное обновление
		return c.JSON(http.StatusOK, map[string]string{"message": "Password updated successfully"})
	}
}

func GetUserInfoByToken(db *pgx.Conn) echo.HandlerFunc {
	fmt.Println("/home/og/projects/Auth/internal/repository/user.go", "GetUserInfoByToken LOADED")
	return func(c echo.Context) error {
		user, err := helpers.GetUserByToken(c, db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}
		var userResponse = models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		}

		return c.JSON(http.StatusOK, userResponse)
	}
}

func ResetPassword(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := helpers.GetUserByToken(c, db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}
		err = helpers.SendMail("sasha_zakirov_2014@mail.ru", user.Email, "password reset", "https://google.com")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error with sending email"})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Password reset mail sent"})
	}
}
