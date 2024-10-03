package routes

import (
	"auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Echo, dbConn *pgx.Conn, rg *echo.Group) {
	e.POST("/auth/create-user", repository.RegisterUser(dbConn))
	e.POST("/auth/login", repository.LoginUser(dbConn))
	rg.POST("/change-password", repository.ChangePassword(dbConn))
	//TODO: reset-password (need mail implementation)
	//TODO: GET  /api/user инфа о пользователе по токену}
}
