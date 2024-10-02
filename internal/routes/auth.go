package routes

import (
	"auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthRoutes(e *echo.Echo, dbConn *pgx.Conn) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/auth/create-user", repository.RegisterUser(dbConn))
}
