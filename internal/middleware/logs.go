package middleware

import (
	"auth/internal/helpers"
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

func AdminMiddleware(dbConn *pgx.Conn) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
				token = token[7:] // Убираем "Bearer " из заголовка
			}

			claims := &helpers.Claims{}
			tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(helpers.JWTSalt), nil
			})

			if err != nil || !tkn.Valid {
				fmt.Println(err)
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			var role int
			err = dbConn.QueryRow(context.Background(), "SELECT role FROM users WHERE id=$1", claims.UserId).Scan(&role)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
			}

			if role != 1 {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
			}

			return next(c)
		}
	}
}

func LogRequestMiddleware(dbConn *pgx.Conn) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("user_id")
			err := next(c)
			start := time.Now()
			method := c.Request().Method
			path := c.Request().URL.Path
			statusCode := c.Response().Status

			if !strings.HasPrefix(path, "/api") {
				if userID != nil {
					_, dbErr := dbConn.Exec(context.Background(), repository.AddLogsQuery, userID, method, path, statusCode, start)
					if dbErr != nil {
						fmt.Printf("Failed to log request: %v\n", dbErr)
					}
				} else {
					fmt.Println("No user ID found in context")
				}
			}
			return err
		}
	}
}
