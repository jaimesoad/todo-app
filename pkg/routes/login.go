package routes

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
	"todolist/pkg/global"
	"todolist/pkg/model"
	"todolist/pkg/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func PostLogin(c echo.Context) error {
	username := c.FormValue("username")
	passwd := c.FormValue("passwd")

	if username == "" {
		return util.LoginWithMessage("Username can't be empty", c)
	}
	if passwd == "" {
		return util.LoginWithMessage("Password can't be empty", c)
	}

	user, err := global.Q.GetUserData(global.Ctx, username)
	if err != nil {
		return util.LoginWithMessage(fmt.Sprintf("User: \"%s\" does not exist", username), c)
	}

	hash := sha256.Sum256([]byte(passwd + user.Salt))

	if string(hash[:]) != string(user.Passwd) {
		return util.LoginWithMessage("Incorrect password", c)
	}

	claims := &model.JwtCustomClaims{
		Username: username,
		Id:       user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(global.SecretKey)
	if err != nil {
		return util.LoginWithMessage(err.Error(), c)
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    t,
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
	}

	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/home")
}
