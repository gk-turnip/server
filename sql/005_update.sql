
-- for version 0.2.11

-- holds user preferences

create table user_prefs
	(
	user_id int,
	pref_name varchar(20),
	pref_value text,
	primary key(user_id, pref_name),
	foreign key(user_id) references users(id)
	);

