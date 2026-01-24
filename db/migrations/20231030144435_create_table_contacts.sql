-- +goose Up
create table contacts
(
    id         char(36)     not null,
    first_name varchar(100) not null,
    last_name  varchar(100) null,
    email      varchar(100) null,
    phone      varchar(100) null,
    user_id    char(36)     not null,
    created_at bigint       not null,
    updated_at bigint       not null,
    primary key (id),
    foreign key fk_contacts_user_id (user_id) references users (id)
) engine = innodb;

-- +goose Down
drop table contacts;
