create table users (
	id bigserial primary key,
	uid uuid not null unique default gen_random_uuid(),
	name text not null,
	email text not null unique,
	password_hash text not null,
	created_at timestamptz default now(),
	updated_at timestamptz default now()
)