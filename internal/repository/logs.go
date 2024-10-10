package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func GetLogs(dbConn *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		rows, err := dbConn.Query(context.Background(), GetLogsQuery)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch logs"})
		}
		defer rows.Close()

		var logs []map[string]interface{}
		for rows.Next() {
			var id, userID, statusCode int
			var method, path string
			var createdAt time.Time
			err := rows.Scan(&id, &userID, &method, &path, &statusCode, &createdAt)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to parse logs"})
			}
			logs = append(logs, map[string]interface{}{
				"id":          id,
				"user_id":     userID,
				"method":      method,
				"path":        path,
				"status_code": statusCode,
				"created_at":  createdAt,
			})
		}

		return c.JSON(http.StatusOK, logs)
	}
}
