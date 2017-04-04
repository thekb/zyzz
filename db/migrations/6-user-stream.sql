-- +migrate Up
alter table users
    add column language text not null default '';

alter table stream
    add column language text not null default '',
    add column match_source text not null default '';

-- +migrate Down
alter table usres
    drop column language;

alter table stream
    drop column language,
    drop column match_source;