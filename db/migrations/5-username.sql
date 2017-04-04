-- +migrate Up
alter table users
add column username text not null default '' unique;

-- +migrate Down
alter table usres
drop column username;