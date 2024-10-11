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

		if req.Password != req.PasswordConfirm {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Passwords don't match"})
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
		var userId int
		var role int
		err = db.QueryRow(context.Background(), GetUserData, req.Email).Scan(&userId, &hashedPassword, &role)
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

		token, err := helpers.GenerateJWT(req.Email, userId, role)
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
// @Produce     json
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

		var req ChangePasswordRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		var oldPassword string
		var Id int
		var role int
		err = db.QueryRow(context.Background(), GetUserData, user.Email).Scan(&Id, &oldPassword, &role)
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid old password"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
		}

		if helpers.CheckPassword(req.NewPassword, oldPassword) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "New password should not repeat old password"})
		}

		newHashedPassword, err := helpers.HashPassword(req.NewPassword)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash new password"})
		}

		_, err = db.Exec(context.Background(), UpdatePasswordQuery, newHashedPassword, user.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
		}

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

// @Summary      Сброс пароля
// @Security ApiKeyAuth
// @Description  Сброс пароля происходит через токен в запросе (токен генерируется в ссылке, которая отправляется на почту)
// @Tags         users
// @ID reset-password
// @Accept       json
// @Produce      json
// @Param input body ResetPasswordRequest true "New password with confirmation and reset token"
// @Success 200 {string} map[string]string
// @Failure 400 {string} map[string]string
// @Failure 401 {string} map[string]string
// @Failure 500 {string} map[string]string
// @Router /api/user/reset-password [put]
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

// @Summary      Отправка письма со сбросом пароля на почту
// @Security ApiKeyAuth
// @Description Почту пользователя получаем через токен
// @Tags         users
// @ID reset-password-link
// @Produce      json
// @Success 200 {string} map[string]string
// @Failure 401 {string} map[string]string
// @Failure 500 {string} map[string]string
// @Router /api/user/send-reset-password-link [get]
func SendResetPasswordLink(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := helpers.GetUserByToken(c, db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}
		token, err := helpers.GenerateJWT(user.Email, user.ID, int(user.Role))
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

func GetUserList(dbConn *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		rows, err := dbConn.Query(context.Background(), GetUserListQuery)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch users"})
		}
		defer rows.Close()

		var users []map[string]interface{}
		for rows.Next() {
			var id int
			var name, email string
			var role int
			err := rows.Scan(&id, &name, &email, &role)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to parse users"})
			}

			roleText := ""
			switch role {
			case 1:
				roleText = "Admin"
			case 2:
				roleText = "Rater"
			case 3:
				roleText = "Watcher"
			default:
				roleText = "Unknown"
			}

			users = append(users, map[string]interface{}{
				"id":    id,
				"name":  name,
				"email": email,
				"role":  roleText,
			})
		}

		// HTML response for htmx
		htmlResponse := ""
		for _, user := range users {
			htmlResponse += fmt.Sprintf(`
                <tr>
                    <td>%d</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`, user["id"], user["name"], user["email"], user["role"])
		}

		return c.HTML(http.StatusOK, htmlResponse)
	}
}
