package global

import (
	"context"
	"text/template"
	sqlc "todolist/pkg/db"
	"todolist/pkg/model"
)

var Ctx context.Context
var Q *sqlc.Queries
var SecretKey []byte
var Templs = &model.Templates{
	Templates: template.Must(template.ParseFiles("template/todoItem.html", "template/todoHome.html")),
}
