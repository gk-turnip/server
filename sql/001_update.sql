
create sequence users_seq
	increment by 1
	minvalue 1
	maxvalue 9223372036854775807
	start with 1;

create table users
	(
	id bigint,
	user_name varchar(20),
	password_hash varchar(128),
	password_salt varchar(10),
	email varchar(80),
	primary key(id)
	);

create unique index users_2 on users (user_name);

