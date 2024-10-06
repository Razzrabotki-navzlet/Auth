package routes

import (
	"github.com/labstack/echo/v4"
)

func AppRoutes(e *echo.Echo, rg *echo.Group) {
	e.GET("/login", ServeLoginPage)
	e.GET("/main", ServeMainPage)
	e.GET("/password-change", ServePasswordChangePage)
}

func ServeLoginPage(c echo.Context) error {
	return c.File("public/html/login.html")
}

func ServeMainPage(c echo.Context) error {
	return c.File("public/html/main.html")
}

func ServePasswordChangePage(c echo.Context) error {
	return c.File("public/html/password-change.html")
}
