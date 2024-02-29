package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"text/template"
	"todolist/pkg/global"
	"todolist/pkg/model"
	"todolist/pkg/routes"
	"todolist/pkg/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

//go:embed sql/schema.sql
var schema string

func main() {
	godotenv.Load()
	global.Ctx, global.Q = util.NewPgx(schema)

	global.SecretKey = []byte(os.Getenv("SECRET"))

	e := echo.New()

	e.Renderer = &model.Templates{
		Templates: template.Must(template.ParseGlob("template/*.html")),
	}

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.JwtCustomClaims)
		},
		SigningKey:  global.SecretKey,
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.Redirect(http.StatusFound, "/login")
		},
	}

	restricted := e.Group("", echojwt.WithConfig(jwtConfig))

	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.Render(http.StatusOK, "404", nil)
	})

	e.Static("/static", "static")

	e.GET("", func(c echo.Context) error {
		//TODO: check if user is logged in, if true, redirect to home, else to login
		return c.Redirect(http.StatusFound, "/home")
	}, echojwt.WithConfig(jwtConfig))

	// Unrestricted routes
	e.GET("/login", routes.GetLogin)
	e.POST("/login", routes.PostLogin)
	e.GET("/signup", routes.GetRegister)
	e.POST("/signup", routes.PostRegister)

	restricted.GET("/logout", routes.GetLogout)
	restricted.GET("/home", routes.GetHome)
	restricted.POST("/todo", routes.PostTodo)
	restricted.DELETE("/todo/:id", routes.DeleteTodoById)
	restricted.PUT("/todo/:id", routes.ToggleTodoById)

	log.Fatal(e.Start(":3000"))
}
