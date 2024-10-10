package middleware

import (
	"auth/internal/repository"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

func LogRequestMiddleware(dbConn *pgx.Conn) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("user_id")
			err := next(c)
			start := time.Now()
			method := c.Request().Method
			path := c.Request().URL.Path
			statusCode := c.Response().Status

			if strings.HasPrefix(path, "/api") && path != "/api/user/current" {
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
