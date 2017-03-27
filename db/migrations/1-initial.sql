-- +migrate Up
create table stream_server(
    id integer primary key autoincrement,
    short_id text not null,
    name text not null,
    hostname text not null,
    internal_ip text not null,
    external_ip text not null
);

create table users(
    id integer primary key autoincrement,
    short_id text not null,
    name text not null,
    description not null,
    created_at datetime not null default current_timestamp,
    published int not null default 0,
    subscribed int not null default 0
);


-- +migrate Down
drop table users;
drop table stream_server;