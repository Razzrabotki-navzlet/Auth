package routes

import (
	localMiddleware "auth/internal/middleware"
	"auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func AdminRoutes(dbConn *pgx.Conn, rg *echo.Group) {
	rg.Use(localMiddleware.AdminMiddleware(dbConn))
	rg.GET("/admin/logs", repository.GetLogs(dbConn))
	rg.GET("/admin/users", repository.GetUserList(dbConn))
	rg.PUT("/admin/change-role", repository.ChangeUserRole(dbConn))
}
