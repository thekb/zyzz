-- +migrate Up
create table event(
    id integer primary key autoincrement,
    short_id text not null,
    name text not null,
    description not null,
    created_at datetime not null default current_timestamp,
    starttime datetime,
    endtime datetime,
    running_nown bool
);

create table event_stream(
    id integer primary key autoincrement,
    short_id text not null,
    event_id integer not null,
    stream_id integer not null,
    foreign key (event_id) references event(id),
    foreign key (stream_id) references stream(id)
);
-- +migrate Down
drop table event;
drop table event_stream;
