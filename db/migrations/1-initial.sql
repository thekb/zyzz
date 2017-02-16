-- +migrate Up
create table stream_server(
    id integer primary key autoincrement,
    short_id text not null,
    name text not null,
    hostname text not null,
    internal_ip text not null,
    external_ip text not null
);

create table user(
    id integer primary key autoincrement,
    short_id text not null,
    name text not null,
    description not null,
    created_at datetime not null default current_timestamp,
    published int not null default 0,
    subscribed int not null default 0
);

create table stream(
    id integer primary key autoincrement,
    short_id text not null,
    name text not null,
    description text not null,
    started_at datetime not null default current_timestamp,
    ended_at datetime not null default current_timestamp,
    status int not null default 0,
    subscriber_count int not null default 0,
    publish_url text not null,
    subscribe_url text not null,
    stream_server_id int not null,
    creator_id int not null,
    transport_url text not null,
    foreign key(creator_id) references user(id),
    foreign key(stream_server_id) references stream_servers(id)
);

-- +migrate Down
drop table stream;
drop table user;
drop table stream_server;