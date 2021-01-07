package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// Template html template type
type Template struct {
	templates *template.Template
}

// Render base render
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
