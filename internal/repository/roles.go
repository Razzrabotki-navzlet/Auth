package repository

import (
	"auth/internal/helpers"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetUserRole(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := helpers.GetUserByToken(c, db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}
		return c.JSON(http.StatusOK, user.Role)
	}
}

func UpdateUserRole(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := helpers.GetUserByToken(c, db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}
		var req ChangeRoleRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		// изменять роли могут только админы
		//if user.Role == models.AdminRole {
		_, err = db.Exec(context.Background(), UpdateRoleQuery, req.NewRole, req.UserId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update role"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Role changed successfully"})
		//} else {
		//	return c.JSON(http.StatusBadRequest, map[string]string{"error": "You don't have permission to change role"})
		//}

	}
}
