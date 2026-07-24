create table Events (
	id SERIAL primary key,
	description VARCHAR,
	tg_userId VARCHAR not null,
	date_create TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    date_end TIMESTAMPTZ not null,
	status BOOL
)