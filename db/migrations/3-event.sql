-- +migrate Up
create table event(
    id serial primary key,
    short_id text not null unique,
    name text not null,
    description text not null,
    created_at timestamp not null default current_timestamp,
    starttime timestamp,
    endtime timestamp,
    running_now int,
    matchid int,
    matchurl text
);

create table stream(
    id serial primary key,
    short_id text not null unique,
    name text not null,
    description text not null,
    started_at timestamp not null default current_timestamp,
    ended_at timestamp not null default current_timestamp,
    status int not null default 0,
    subscriber_count int not null default 0,
    publish_url text not null,
    subscribe_url text not null,
    stream_server_id int not null,
    creator_id int not null,
    transport_url text not null,
    event_id text not null,
    foreign key(creator_id) references users(id),
    foreign key(stream_server_id) references stream_server(id),
    foreign key(event_id) references event(short_id)
);
-- +migrate Down
drop table event;
drop table stream;