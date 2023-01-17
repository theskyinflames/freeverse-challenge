CREATE TABLE if not exists products (
	id uuid NOT NULL,
	name VARCHAR(100) not null,
	price  numeric(15,6) null,
	available BOOLEAN default false,
	PRIMARY KEY (id)
);