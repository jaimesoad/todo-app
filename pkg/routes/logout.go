package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetLogout(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/login")
	}

	cookie.Expires = time.Now()

	c.SetCookie(cookie)

	return c.Redirect(http.StatusPermanentRedirect, "/login")
}
