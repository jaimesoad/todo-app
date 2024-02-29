package routes

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	sqlc "todolist/pkg/db"
	"todolist/pkg/global"
	"todolist/pkg/util"

	"github.com/labstack/echo/v4"
)

func PostTodo(c echo.Context) error {
	var todo sqlc.NewTodoParams

	user := util.GetSessionUser(c)

	todo.UserID = user.Id
	todo.Content = c.FormValue("content")

	if todo.Content == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	newTodo, err := global.Q.NewTodo(global.Ctx, todo)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var buf bytes.Buffer

	err = global.Templs.Execute(&buf, "todoItem", newTodo)
	if err != nil {
		fmt.Println(err.Error())
	}

	return c.HTML(http.StatusOK, buf.String())
}

func DeleteTodoById(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user := util.GetSessionUser(c)

	delTodo := sqlc.DeleteTodoByIdParams{
		UserID: user.Id,
		ID:     int32(id),
	}

	err = global.Q.DeleteTodoById(global.Ctx, delTodo)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.HTML(http.StatusOK, "")
}

func ToggleTodoById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user := util.GetSessionUser(c)

	toggleTodo := sqlc.ToggleUserTodoParams{
		UserID: user.Id,
		ID:     int32(id),
	}

	done, err := global.Q.ToggleUserTodo(global.Ctx, toggleTodo)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	content := "mark as done"

	if done {
		content = "mark as undone"
	}

	return c.HTML(http.StatusOK, content)
}