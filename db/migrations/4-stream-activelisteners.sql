-- +migrate Up
alter table stream
add column active_listeners int not null default 0;

-- +migrate Down
alter table stream
drop column active_listeners;