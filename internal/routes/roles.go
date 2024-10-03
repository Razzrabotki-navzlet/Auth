package routes

import (
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func RolesRoutes(dbConn *pgx.Conn, rg *echo.Group) {
	// инфа о роли пользователя по почте? или по айди?
	// rg.GET("/user-role", repositore.GetUserRole)

	// инфа о роли пользователя по почте? или по айди?
	// админы могут изменять роль модераторов и юзеров, просматривать логи
	// модераторы могут изменять роль юзеров
	// rg.PUT("/change-user-role", repositore.ChangeUserRole)
}
