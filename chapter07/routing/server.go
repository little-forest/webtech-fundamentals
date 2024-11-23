package main

import (
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	port := getPortNumber()

	e := echo.New()

	// Setup middlewares
	e.Use(middleware.Logger())

	// Prepare static contents
	e.Static("/", "static")

	// Rewrite url
	e.Pre(middleware.Rewrite(map[string]string{
		"/page*": "/index.html",
	}))

	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
