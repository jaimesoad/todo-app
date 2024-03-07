create table Users (
    id       serial      primary key,
    username varchar(30) unique not null,
    passwd   bytea       not null,
    salt     varchar(10) not null
);

create table Todos (
    id      serial  primary key,
    content text    not null,
    done    boolean not null default false,
    user_id integer not null,
    foreign key (user_id) references Users (id)
);
create index user_todo on Todos (id, user_id);

create view Users_View as
select id, username
from Users;