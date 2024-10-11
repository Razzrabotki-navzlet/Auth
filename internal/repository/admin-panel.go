package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetLogs(dbConn *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		rows, err := dbConn.Query(context.Background(), GetLogsQuery)
		if err != nil {
			return c.HTML(http.StatusInternalServerError, "<p>Failed to fetch logs</p>")
		}
		defer rows.Close()

		var htmlOutput strings.Builder
		for rows.Next() {
			var id, userID, statusCode int
			var method, path string
			var createdAt time.Time
			err := rows.Scan(&id, &userID, &method, &path, &statusCode, &createdAt)
			if err != nil {
				return c.HTML(http.StatusInternalServerError, "<p>Failed to parse logs</p>")
			}

			htmlOutput.WriteString("<tr>")
			htmlOutput.WriteString(fmt.Sprintf("<td>%d</td>", id))
			htmlOutput.WriteString(fmt.Sprintf("<td>%d</td>", userID))
			htmlOutput.WriteString(fmt.Sprintf("<td>%s</td>", method))
			htmlOutput.WriteString(fmt.Sprintf("<td>%s</td>", path))
			htmlOutput.WriteString(fmt.Sprintf("<td>%d</td>", statusCode))
			htmlOutput.WriteString(fmt.Sprintf("<td>%s</td>", createdAt.Format("2006-01-02 15:04:05")))
			htmlOutput.WriteString("</tr>")
		}

		if htmlOutput.Len() == 0 {
			htmlOutput.WriteString("<tr><td colspan='6'>Нет данных для отображения</td></tr>")
		}

		return c.HTML(http.StatusOK, htmlOutput.String())
	}
}

func ChangeUserRole(dbConn *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			UserID  string `json:"user_id"`
			NewRole string `json:"new_role"`
		}

		if err := c.Bind(&req); err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
		}

		userID, err := strconv.Atoi(req.UserID)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id format"})
		}

		newRole, err := strconv.Atoi(req.NewRole)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid new_role format"})
		}

		if newRole < 1 || newRole > 3 {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid role. Choose between 1 (Admin), 2 (Rater), 3 (Watcher)"})
		}
		result, err := dbConn.Exec(context.Background(), ChangeRoleQuery, newRole, userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update user role"})
		}

		rowsAffected := result.RowsAffected()
		if rowsAffected == 0 {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
		}

		return c.JSON(http.StatusOK, echo.Map{"message": "User role updated successfully"})
	}
}
