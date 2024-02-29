package routes

import (
	"bytes"
	"net/http"
	"todolist/pkg/global"
	"todolist/pkg/model"
	"todolist/pkg/util"

	"github.com/labstack/echo/v4"
)

func GetHome(c echo.Context) error {
	user := util.GetSessionUser(c)

	todos, err := global.Q.GetUserTodos(global.Ctx, user.Id)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	home := model.Home{
		Username: user.Username,
		Id:       user.Id,
	}

	for _, todo := range todos {
		var todoHTML bytes.Buffer

		global.Templs.Execute(&todoHTML, "todoHome", todo)
		home.TodoList = append(home.TodoList, todoHTML.String())
	}

	return c.Render(http.StatusOK, "home", home)
}
