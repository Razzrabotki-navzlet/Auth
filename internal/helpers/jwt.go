package helpers

import (
	"auth/internal/models"
	"auth/internal/repository"
	"context"
	"fmt"
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
	Email  string `json:"email"`
	UserId int    `json:"userId"`
	Role   int    `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(email string, userId int, role int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email:  email,
		UserId: userId,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}

	fmt.Println("JWT token generated successfully for email:", email)
	return tokenString, nil
}

func GetUserByToken(c echo.Context, db *pgx.Conn) (models.User, error) {
	claims := &Claims{}
	token := c.Request().Header.Get("Authorization")
	var user models.User

	if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	} else {
		return user, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
	}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSalt), nil
	})

	if err != nil || !tkn.Valid {
		return user, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	email := claims.Email

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

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		} else {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token format"})
		}

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSalt), nil
		})

		if err != nil || !tkn.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		userID := claims.UserId
		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User ID not found in token"})
		}

		c.Set("user_id", userID)

		return next(c)
	}
}

func GenerateConfirmationJWT(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}

	return tokenString, nil
}

func ConfirmEmail(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSalt), nil
		})

		if err != nil || !tkn.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
		}

		_, err = db.Exec(context.Background(), repository.UpdateVerifyQuery, claims.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to activate user"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Email confirmed successfully"})
	}
}
