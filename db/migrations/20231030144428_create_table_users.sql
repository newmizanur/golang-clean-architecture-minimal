-- +goose Up
create table users
(
    id         char(36)     not null,
    name       varchar(100) not null,
    password   varchar(100) not null,
    created_at bigint       not null,
    updated_at bigint       not null,
    primary key (id)
) engine = InnoDB;

-- +goose Down
drop table users;
