package routes

import (
	"auth/internal/helpers"
	"auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Echo, dbConn *pgx.Conn, rg *echo.Group) {
	e.POST("/auth/create-user", repository.RegisterUser(dbConn))
	e.POST("/auth/login", repository.LoginUser(dbConn))
	e.GET("/user/confirm", helpers.ConfirmEmail(dbConn))
	rg.PUT("/user/change-password", repository.ChangePassword(dbConn))
	rg.GET("/user/current", repository.GetUserInfoByToken(dbConn))
	rg.PUT("/user/reset-password", repository.ResetPassword(dbConn))
	rg.GET("/user/send-reset-password-link", repository.SendResetPasswordLink(dbConn))
}
