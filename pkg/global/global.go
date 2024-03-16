package global

import (
	"text/template"
	"todolist/pkg/model"
)

var SecretKey []byte
var Templs = &model.Templates{
	Templates: template.Must(template.ParseFiles("template/todoItem.html", "template/todoHome.html")),
}
