package mid

import (
	"context"
	sqlc "todolist/pkg/db"

	"github.com/labstack/echo/v4"
)

type DBConfig struct {
	Conn *sqlc.Queries
	Ctx context.Context
}

func WithDBConfig(config DBConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			c.Set("db", config)
			return next(c)
		}
	}
}