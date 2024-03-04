package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetLogout(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	cookie.Value = ""
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	cookie.Path = "/"

	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/login")
}
