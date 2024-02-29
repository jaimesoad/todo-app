package model

import (
	"text/template"
	"io"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Login struct {
	Message string `json:"message"`
}

type SessionUser struct {
	Username string `json:"username"`
	Id       int32  `json:"id"`
}

type JwtCustomClaims struct {
	Username string `json:"username"`
	Id       int32  `json:"id"`
	jwt.RegisteredClaims
}

type Home struct {
	Username string
	Id       int32
	TodoList []string
}

type JSON map[string]interface{}

type Templates struct {
	Templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name+".html", data)
}

func (t *Templates) Execute(w io.Writer, name string, data any) error {
	return t.Templates.ExecuteTemplate(w, name+".html", data)
}
