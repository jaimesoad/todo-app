-- name: NewUser :exec
insert into Users(username, passwd, salt)
values($1, $2, $3);

-- name: GetUserData :one
select * from Users
where username = $1;

-- name: NewTodo :one
insert into Todos(content, user_id)
values($1, $2)
returning *;

-- name: DeleteTodoById :exec
delete from Todos
where id = $1 and user_id = $2;

-- name: GetUserTodos :many
select id, content, done
from Todos
where user_id = $1
order by id desc;

-- name: ToggleUserTodo :one
update Todos
set done = not done
where id = $1 and user_id = $2
returning done;