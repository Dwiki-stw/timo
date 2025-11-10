create table photos (
	id bigserial primary key,
	journal_id bigint not null references journals(id),
	url text not null,
	created_at timestamptz default now()
)