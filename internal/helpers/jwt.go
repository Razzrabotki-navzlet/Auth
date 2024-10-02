package helpers

import (
	"auth/internal/models"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

const JWTSalt = "baklajan"
const GetInfoByTokenQuery = `SELECT id, name, email, role, password FROM users WHERE email = $1`

var jwtKey = []byte(JWTSalt)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Токен будет действителен 24 часа
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserByToken(c echo.Context, db *pgx.Conn) (models.User, error) {
	claims := &Claims{}
	token := c.Request().Header.Get("Authorization")
	var user models.User

	// Проверка формата токена
	if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	} else {
		return user, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
	}

	// Валидация и декодирование токена
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("baklajan"), nil
	})

	if err != nil || !tkn.Valid {
		return user, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	// Получаем userID из claims
	email := claims.Email

	// Запрос в базу данных для получения информации о пользователе
	err = db.QueryRow(context.Background(), GetInfoByTokenQuery, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Password,
	)

	if err != nil {
		return user, echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user info")
	}

	return user, nil
}
