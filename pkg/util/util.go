package util

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
	sqlc "todolist/pkg/db"
	mid "todolist/pkg/middleware"
	"todolist/pkg/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewPgx(databse string) (context.Context, *sqlc.Queries) {

	connstr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DBHOST"),
		os.Getenv("POSTGRES_USER"),
	)

	db, err := pgx.Connect(context.Background(), connstr)
	if err != nil {
		panic(err.Error())
	}

	db.Exec(context.Background(), databse)

	return context.Background(), sqlc.New(db)
}

func NewSalt() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func LoginWithMessage(message string, c echo.Context) error {
	return c.Render(http.StatusOK, "login", model.Login{Message: message})
}

func SignupWithMessage(message string, c echo.Context) error {
	return c.Render(http.StatusOK, "register", model.Login{Message: message})
}

func GetSessionUser(c echo.Context) model.SessionUser {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)

	return model.SessionUser{
		Username: claims.Username,
		Id:       claims.Id,
	}
}

func GetDBSession(c echo.Context) (*sqlc.Queries, context.Context) {
	db := c.Get("db").(mid.DBConfig)

	return db.Conn, db.Ctx
}