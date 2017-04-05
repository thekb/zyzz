-- +migrate Up
alter table users
add column username text not null default md5(random()::text) unique;

-- +migrate Down
alter table usres
drop column username;