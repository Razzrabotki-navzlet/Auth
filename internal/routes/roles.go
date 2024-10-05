package routes

import (
	"auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func RolesRoutes(dbConn *pgx.Conn, rg *echo.Group) {

	rg.GET("/user-role", repository.GetUserRole(dbConn))
	rg.PUT("/update-user-role", repository.UpdateUserRole(dbConn))

}
