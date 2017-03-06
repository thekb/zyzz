-- +migrate Up
alter table user add column email text;
alter table user add column nickname text;
alter table user add column avatarurl text;
alter table user add column fbid text;
alter table user add column access_token text;

-- +migrate Down
alter table user drop column email;
alter table user drop column nickname;
alter table user drop column avatarurl;
alter table user drop column fbid;
alter table user drop column access_token;