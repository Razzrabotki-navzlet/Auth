package main

import (
	_ "auth/docs"
	"auth/internal/helpers"
	"auth/internal/routes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
)

// @title           API
// @version         123.123
// @description     This is a sample API documentation using Swagger.

// @host      localhost:7070
// @BasePath  /

//@securityDefinitions.apikey ApiKeyAuth
//@in header
//@name Authorization

func main() {
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:7070", "http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderContentType},
	}))
	rg := e.Group("/api")
	rg.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("baklajan"), // Указываем ключ подписи
	}))
	rg.Use(helpers.JWTMiddleware)

	dbUrl := fmt.Sprintf("postgres://admin:root@localhost:5432/auth")
	//dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
	//	os.Getenv("POSTGRES_USER"),
	//	os.Getenv("POSTGRES_PASSWORD"),
	//	os.Getenv("POSTGRES_HOST"),
	//	os.Getenv("POSTGRES_PORT"),
	//	os.Getenv("POSTGRES_DB"),
	//)
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	routes.AuthRoutes(e, conn, rg)
	routes.RolesRoutes(conn, rg)
	routes.AppRoutes(e, rg)
	e.Logger.Fatal(e.Start(":7070"))
	defer conn.Close(context.Background())
}
