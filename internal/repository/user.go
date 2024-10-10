package repository

import (
	"auth/internal/helpers"
	"auth/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary      Регистрация
// @Description  Создание пользователя
// @Tags         users
// @ID register-user
// @Accept       json
// @Produce      json
// @Param input body RegisterRequest true "user info"
// @Success 200 {string} map[string]string
// @Failure 400 {string} map[string]string
// @Failure 500 {string} map[string]string
// @Router /auth/create-user [post]
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

		token, err := helpers.GenerateConfirmationJWT(req.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate confirmation token"})
		}

		confirmationLink := fmt.Sprintf("http://localhost:7070/user/confirm?token=%s", token)

		err = helpers.SendMail([]string{req.Email}, "Email Confirmation", "Please confirm your email by clicking the following link: "+confirmationLink)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send confirmation email"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User created successfully. Please check your email to confirm your account."})
	}
}

// @Summary      Логин
// @Description  Получение JWT
// @Tags         users
// @ID login-user
// @Accept       json
// @Produce      json
// @Param input body LoginRequest true "user info"
// @Success 200 {string} map[string]string
// @Failure 400 {string} map[string]string
// @Failure 401 {string} map[string]string
// @Failure 500 {string} map[string]string
// @Router /auth/login [post]
func LoginUser(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			fmt.Println("Error binding request:", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}
		isVerified := false
		err := db.QueryRow(context.Background(), CheckVerify, req.Email).Scan(&isVerified)
		if !isVerified {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "You should verify your account (check mail)"})
		}
		var hashedPassword string
		err = db.QueryRow(context.Background(), CheckPasswordQuery, req.Email).Scan(&hashedPassword)
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

// @Summary      Смена пароля
// @Security ApiKeyAuth
// @Tags         users
// @ID password-change
// @Accept       json
// @Producea      json
// @Param input body ChangePasswordRequest true "Ols and new passwords"
// @Success 200 {string} map[string]string
// @Failure 400 {string} map[string]string
// @Failure 401 {string} map[string]string
// @Failure 500 {string} map[string]string
// @Router /api/user/change-password [put]
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

// @Summary      Информация о пользователе
// @Security ApiKeyAuth
// @Description  Получение информации о текущем пользователе по токену
// @Tags         users
// @ID info-user
// @Produce      json
// @Success 200 {string} map[string]string
// @Failure 400 {string} map[string]string
// @Failure 401 {string} map[string]string
// @Router /api/user/current [get]
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

		var req ResetPasswordRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		if req.NewPassword != req.ConfirmPassword {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Passwords do not match"})
		}

		claims := &helpers.Claims{}
		fmt.Println("Received token:", req.Token)
		tkn, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(helpers.JWTSalt), nil
		})

		if err != nil || !tkn.Valid {
			fmt.Println("Token:", tkn)
			fmt.Println("ERROR:", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
		}

		hashedPassword, err := helpers.HashPassword(req.NewPassword)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}

		_, err = db.Exec(context.Background(), "UPDATE users SET password = $1 WHERE email = $2", hashedPassword, claims.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
	}
}

func SendResetPasswordLink(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := helpers.GetUserByToken(c, db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}
		token, err := helpers.GenerateJWT(user.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate reset token"})
		}

		resetLink := fmt.Sprintf("http://localhost:7070/password-reset?token=%s", token)
		body := fmt.Sprintf("Для сброса пароля перейдите по ссылке: %s", resetLink)
		err = helpers.SendMail([]string{user.Email}, "Сброс пароля", body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Reset password link sent to your email"})
	}
}
