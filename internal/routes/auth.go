package routes

import (
	"auth/internal/helpers"
	"auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Echo, dbConn *pgx.Conn, rg *echo.Group) {
	rg.Use(helpers.JWTMiddleware)
	e.GET("/login", ServeLoginPage)
	e.GET("/main", ServeMainPage)
	e.POST("/auth/create-user", repository.RegisterUser(dbConn))
	e.POST("/auth/login", repository.LoginUser(dbConn))
	rg.PUT("/change-password", repository.ChangePassword(dbConn))
	rg.GET("/user/current", repository.GetUserInfoByToken(dbConn))
	//TODO настроить сервер для возможности отправки писем
	rg.GET("/user/reset-password", repository.ResetPassword(dbConn))
}

func ServeLoginPage(c echo.Context) error {
	return c.File("public/html/login.html")
}

func ServeMainPage(c echo.Context) error {
	return c.File("public/html/main.html")
}
