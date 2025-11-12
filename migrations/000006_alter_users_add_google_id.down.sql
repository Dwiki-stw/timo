alter table users
alter column password_hash set not null;

alter table users
drop column google_id;