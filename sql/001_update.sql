
-- updated on version 0.2.15

-- users

create sequence users_seq
	increment by 1
	minvalue 1
	maxvalue 2147483647
	start with 1;

create table users
	(
	id int not null,
	user_name varchar(20) not null,
	password_hash varchar(128) not null,
	password_salt varchar(10) not null,
	email varchar(80) not null,
	account_creation_date timestamp not null,
	last_login_date timestamp not null,
	primary key(id)
	);

create unique index users_2 on users (user_name);

-- user preferences

create table user_prefs
	(
	user_id int not null,
	pref_name varchar(20) not null,
	pref_value text not null,
	primary key(user_id, pref_name),
	foreign key(user_id) references users(id)
	);

-- chat archives

create sequence chat_archives_seq
	increment by 1
	minvalue 1
	maxvalue 2147483647
	start with 1;

create table chat_archives
	(
	id int not null,
	user_id int not null,
	message_creation_date timestamp not null,
	chat_message text not null,
	primary key(id),
	foreign key(user_id) references users(id)
	);

-- pods

drop table pods;

create table pods
	(
	id int not null,
	pod_title text not null,
	editable boolean not null,
	primary key(id)
	);

