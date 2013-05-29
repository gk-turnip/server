
-- for version 0.1.27

create sequence chat_archives_seq
	increment by 1
	minvalue 1
	maxvalue 2147483647
	start with 1;

create table chat_archives
	(
	id int,
	user_id int,
	message_creation_date timestamp,
	chat_message text,
	primary key(id),
	foreign key(user_id) references users(id)
	);

