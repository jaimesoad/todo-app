package main

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
	sqlc "todolist/pkg/db"
	"todolist/pkg/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Login struct {
	Message string `json:"message"`
}

type SessionUser struct {
	Username string `json:"username"`
	Id       int32  `json:"id"`
}

type jwtCustomClaims struct {
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

//go:embed sql/schema.sql
var schema string

var secretKey []byte

func main() {
	godotenv.Load()
	ctx, q := util.NewPgx(schema)

	secretKey = []byte(os.Getenv("SECRET"))

	e := echo.New()

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey:  secretKey,
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.Redirect(http.StatusFound, "/login")
		},
	}

	e.RouteNotFound("/*", func(c echo.Context) error {
		templ := template.Must(template.ParseFiles("template/404.html"))

		return templ.Execute(c.Response(), nil)
	})

	e.Static("/static", "static")

	e.GET("", func(c echo.Context) error {
		//TODO: check if user is logged in, if true, redirect to home, else to login
		return c.Redirect(http.StatusPermanentRedirect, "/home")
	}, echojwt.WithConfig(jwtConfig))

	e.GET("/login", func(c echo.Context) error {

		templ := template.Must(template.ParseFiles("template/login.html"))

		return templ.Execute(c.Response(), nil)
	})

	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		passwd := c.FormValue("passwd")

		if username == "" {
			return LoginWithMessage("Username can't be empty", c.Response())
		}
		if passwd == "" {
			return LoginWithMessage("Password can't be empty", c.Response())
		}

		user, err := q.GetUserData(ctx, username)
		if err != nil {
			return LoginWithMessage(fmt.Sprintf("User: \"%s\" does not exist", username), c.Response())
		}

		hash := sha256.Sum256([]byte(passwd + user.Salt))

		if string(hash[:]) != string(user.Passwd) {
			return LoginWithMessage("Incorrect password", c.Response())
		}

		claims := &jwtCustomClaims{
			username,
			user.ID,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString(secretKey)
		if err != nil {
			return LoginWithMessage(err.Error(), c.Response())
		}

		cookie := &http.Cookie{
			Name:     "token",
			Value:    t,
			Expires:  time.Now().Add(time.Minute * 30),
			HttpOnly: true,
		}

		c.SetCookie(cookie)

		return c.Redirect(http.StatusFound, "/home")
	})

	e.GET("/signup", func(c echo.Context) error {
		templ := template.Must(template.ParseFiles("template/register.html"))

		return templ.Execute(c.Response(), nil)
	})

	e.POST("/signup", func(c echo.Context) error {
		var credentials sqlc.NewUserParams
		var confirm string

		credentials.Username = strings.ToLower(c.FormValue("username"))
		credentials.Passwd = []byte(c.FormValue("passwd"))
		confirm = c.FormValue("confirm")

		_, err := q.GetUserData(ctx, credentials.Username)
		if err == nil {
			return SignupWithMessage(fmt.Sprintf("User: \"%s\" already exists", credentials.Username), c.Response())
		}

		if len(credentials.Username) == 0 {
			return SignupWithMessage("Username can't be empty", c.Response())
		}

		if len(credentials.Passwd) == 0 {
			return SignupWithMessage("Password can't be empty", c.Response())
		}

		if confirm != string(credentials.Passwd) {
			return SignupWithMessage("Passwords must match", c.Response())
		}

		credentials.Salt = util.NewSalt()

		hashedPasswd := sha256.Sum256(append(credentials.Passwd, []byte(credentials.Salt)...))

		credentials.Passwd = hashedPasswd[:]

		err = q.NewUser(ctx, credentials)
		if err != nil {
			return SignupWithMessage("An error occured. User not created", c.Response())
		}

		return c.Redirect(http.StatusFound, "/login")
	})

	e.GET("/logout", func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusBadRequest, "/login")
		}

		cookie.Expires = time.Now()

		c.SetCookie(cookie)

		return c.Redirect(http.StatusPermanentRedirect, "/login")
	}, echojwt.WithConfig(jwtConfig))

	e.GET("/home", func(c echo.Context) error {
		user := GetSessionUser(c)

		templ := template.Must(template.ParseFiles("template/home.html"))

		todos, err := q.GetUserTodos(ctx, user.Id)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		todoHome := template.Must(template.ParseFiles("template/todoHome.html"))

		home := Home{
			Username: user.Username,
			Id:       user.Id,
		}

		for _, todo := range todos {
			var todoHTML bytes.Buffer

			todoHome.Execute(&todoHTML, todo)
			home.TodoList = append(home.TodoList, todoHTML.String())
		}

		return templ.Execute(c.Response(), home)
	}, echojwt.WithConfig(jwtConfig))

	e.POST("/todo", func(c echo.Context) error {
		var todo sqlc.NewTodoParams

		user := GetSessionUser(c)

		todo.UserID = user.Id
		todo.Content = c.FormValue("content")

		if todo.Content == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		newTodo, err := q.NewTodo(ctx, todo)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		templ := template.Must(template.ParseFiles("template/todoItem.html"))

		var buf bytes.Buffer

		err = templ.Execute(&buf, newTodo)
		if err != nil {
			fmt.Println(err.Error())
		}

		return c.HTML(http.StatusOK, buf.String())
	}, echojwt.WithConfig(jwtConfig))

	e.DELETE("/todo/:id", func(c echo.Context) error {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		user := GetSessionUser(c)

		delTodo := sqlc.DeleteTodoByIdParams{
			UserID: user.Id,
			ID:     int32(id),
		}

		err = q.DeleteTodoById(ctx, delTodo)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.HTML(http.StatusOK, "")
	}, echojwt.WithConfig(jwtConfig))

	e.PUT("/todo/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		user := GetSessionUser(c)

		toggleTodo := sqlc.ToggleUserTodoParams{
			UserID: user.Id,
			ID:     int32(id),
		}

		done, err := q.ToggleUserTodo(ctx, toggleTodo)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		content := "mark as done"

		if done {
			content = "mark as undone"
		}

		return c.HTML(http.StatusOK, content)
	}, echojwt.WithConfig(jwtConfig))

	log.Fatal(e.Start(":3000"))
	fmt.Println("Hello, world!")
}

func LoginWithMessage(message string, wr io.Writer) error {
	templ := template.Must(template.ParseFiles("template/login.html"))

	return templ.Execute(wr, Login{Message: message})
}

func SignupWithMessage(message string, wr io.Writer) error {
	templ := template.Must(template.ParseFiles("template/register.html"))

	return templ.Execute(wr, Login{Message: message})
}

func GetSessionUser(c echo.Context) SessionUser {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)

	return SessionUser{
		Username: claims.Username,
		Id:       claims.Id,
	}
}
