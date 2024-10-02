package routes

import (
	"auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthRoutes(e *echo.Echo, dbConn *pgx.Conn, rg *echo.Group) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/auth/create-user", repository.RegisterUser(dbConn))
	e.POST("/auth/login", repository.LoginUser(dbConn))
	rg.GET("/testPath", func(c echo.Context) error {
		return c.String(http.StatusOK, "TestPath11221")
	})
	rg.POST("/change-password", repository.ChangePassword(dbConn))
	//TODO: change-password
	//TODO: reset-password (need mail implementation)
}
