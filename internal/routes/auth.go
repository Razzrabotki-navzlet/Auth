package routes

import (
	"auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Echo, dbConn *pgx.Conn, rg *echo.Group) {
	e.POST("/auth/create-user", repository.RegisterUser(dbConn))
	e.POST("/auth/login", repository.LoginUser(dbConn))
	rg.PUT("/change-password", repository.ChangePassword(dbConn))
	rg.GET("/user/current", repository.GetUserInfoByToken(dbConn))
	//TODO настроить сервер для возможности отправки писем
	rg.GET("/user/reset-password", repository.ResetPassword(dbConn))
}
