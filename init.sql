-- Скрипт создания таблицы request_logs
create table request_logs
(
    id          serial primary key,
    user_id     integer,
    method      varchar(10),
    path        text,
    status_code integer,
    created_at  timestamp default CURRENT_TIMESTAMP
);

alter table request_logs owner to admin;

-- Скрипт создания таблицы users
create table users
(
    id          serial primary key,
    name        varchar(100) not null,
    email       varchar(100) not null unique,
    role        smallint not null,
    password    varchar(255) not null,
    is_verified boolean default false
);

alter table users owner to admin;
