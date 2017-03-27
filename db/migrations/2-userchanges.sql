-- +migrate Up
alter table users add column email text;
alter table users add column nickname text;
alter table users add column avatarurl text;
alter table users add column fbid text;
alter table users add column access_token text;

-- +migrate Down
alter table users drop column email;
alter table users drop column nickname;
alter table users drop column avatarurl;
alter table users drop column fbid;
alter table users drop column access_token;