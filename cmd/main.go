package main

import (
	"auth/internal/routes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	e := echo.New()
	dbUrl := "postgres://admin:root@localhost:5432/auth"
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	routes.AuthRoutes(e, conn)
	e.Logger.Fatal(e.Start(":7070"))
	defer conn.Close(context.Background())
}
