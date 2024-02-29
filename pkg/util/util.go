package util

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
	sqlc "todolist/pkg/db"

	"github.com/jackc/pgx/v5"
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
