package main

import (
	"auth/internal/routes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	rg := e.Group("/api")
	dbUrl := "postgres://admin:root@localhost:5432/auth"
	rg.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("baklajan"), // Указываем ключ подписи
	}))
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	routes.AuthRoutes(e, conn, rg)
	e.Logger.Fatal(e.Start(":7070"))
	defer conn.Close(context.Background())
}
