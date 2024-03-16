package routes

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	sqlc "todolist/pkg/db"
	"todolist/pkg/util"

	"github.com/labstack/echo/v4"
)

func GetRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register", nil)
}

func PostRegister(c echo.Context) error {
	var credentials sqlc.NewUserParams
	var confirm string

	credentials.Username = strings.ToLower(c.FormValue("username"))
	credentials.Passwd = []byte(c.FormValue("passwd"))
	confirm = c.FormValue("confirm")

	q, ctx := util.GetDBSession(c)

	_, err := q.GetUserData(ctx, credentials.Username)
	if err == nil {
		return util.SignupWithMessage(fmt.Sprintf("User: \"%s\" already exists", credentials.Username), c)
	}

	if len(credentials.Username) == 0 {
		return util.SignupWithMessage("Username can't be empty", c)
	}

	if len(credentials.Passwd) == 0 {
		return util.SignupWithMessage("Password can't be empty", c)
	}

	if confirm != string(credentials.Passwd) {
		return util.SignupWithMessage("Passwords must match", c)
	}

	credentials.Salt = util.NewSalt()

	hashedPasswd := sha256.Sum256(append(credentials.Passwd, []byte(credentials.Salt)...))

	credentials.Passwd = hashedPasswd[:]

	err = q.NewUser(ctx, credentials)
	if err != nil {
		return util.SignupWithMessage("An error occured. User not created", c)
	}

	return c.Redirect(http.StatusFound, "/login")
}
