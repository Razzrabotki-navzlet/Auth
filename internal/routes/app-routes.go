package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"net/http"
)

func AppRoutes(e *echo.Echo, rg *echo.Group) {
	e.GET("/login", RenderPage("login.html"))
	e.GET("/main", RenderPage("main.html"))
	e.GET("/password-change", RenderPage("password-change.html"))
	e.GET("/password-reset", RenderPage("password-reset.html"))
}

func RenderPage(templateName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
		userInput := c.QueryParam("input")
		safeOutput := template.HTMLEscapeString(userInput)
		return c.Render(http.StatusOK, templateName, map[string]interface{}{
			"csrf_token": csrfToken,
			"output":     safeOutput,
		})
	}
}
