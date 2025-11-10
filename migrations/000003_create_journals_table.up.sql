create table journals (
	id bigserial primary key,
	uid uuid not null unique default gen_random_uuid(),
	user_id bigint not null references users(id) on delete cascade,
	title text not null,
	text text not null,
	mood_id bigint not null references moods(id),
	created_at timestamptz default now(),
	updated_at timestamptz default now()
)