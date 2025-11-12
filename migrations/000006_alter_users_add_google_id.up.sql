alter table users
alter column password_hash drop not null;

alter table users
add column google_id text unique;