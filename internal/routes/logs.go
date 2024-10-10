package routes

import (
	"auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func LogRoutes(dbConn *pgx.Conn, rg *echo.Group) {
	rg.GET("/admin/logs", repository.GetLogs(dbConn))
}
