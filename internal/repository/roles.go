package repository

import (
	"auth/internal/helpers"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary      Информация о роли пользователя
// @Security ApiKeyAuth
// @Tags         users
// @ID user-role
// @Produce      json
// @Success 200 {string} map[string]string
// @Failure 401 {string} map[string]string
// @Router /api/user-role [get]
func GetUserRole(db *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := helpers.GetUserByToken(c, db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}
		return c.JSON(http.StatusOK, user.Role)
	}
}

// @Summary      Изменение роли пользователя
// @Security ApiKeyAuth
// @Tags         users
// @ID change-user-role
// @Accept       json
// @Produce      json
// @Param input body ChangeRoleRequest true "user info"
// @Success 200 {string} map[string]string
// @Failure 400 {string} map[string]string
// @Failure 401 {string} map[string]strin
// @Failure 500 {string} map[string]string
// @Router /api/update-user-role [put]
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
		_, err = db.Exec(context.Background(), UpdateRoleQuery, req.NewRole, req.UserId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update role"})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Role changed successfully"})
	}
}
