package routes

import (
	"auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func RolesRoutes(dbConn *pgx.Conn, rg *echo.Group) {
	rg.GET("/user-role", repository.GetUserRole(dbConn))
	rg.PUT("/update-user-role", repository.UpdateUserRole(dbConn))

	// инфа о роли пользователя по почте? или по айди?
	// админы могут изменять роль модераторов и юзеров, просматривать логи
	// модераторы могут изменять роль юзеров
	// rg.PUT("/change-user-role", repositore.ChangeUserRole)
}
