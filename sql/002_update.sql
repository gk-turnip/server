
-- for version 0.1.27

create table context_users
	(
	id int,
	last_position_x smallint,
	last_position_y smallint,
	last_position_z smallint,
	last_pod varchar(4),
	primary key(id),
	foreign key(id) references users(id)
	);

