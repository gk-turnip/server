-- for version 0.2.15

create table attribute_objects
	(
	id integer not null,
	attribute_name varchar(64) not null,
	primary key (id),
	unique (attribute_name)
	);

create table terrain_objects
    (
    id integer not null,
    terrain_object_name varchar(64) not null,
    primary key (id),
	unique (terrain_object_name)
    );

create table pod_tiles
    (
    pod_id integer not null,
    suffix_no integer not null,
    x smallint not null,
    y smallint not null,
	object_id integer not null,
    primary key (pod_id, suffix_no),
    foreign key(pod_id) references pods(id),
	unique (pod_id, x, y)
    );

create table pod_z_tiles
    (
    pod_id integer not null,
    suffix_no integer not null,
    z smallint not null,
	attribute_id smallint not null,
    primary key (pod_id, suffix_no, z),
    foreign key(pod_id, suffix_no) references pod_tiles(pod_id, suffix_no),
	unique (pod_id, z)
    );

create table pod_objects
    (
    pod_id integer not null,
    suffix_no integer not null,
    object_id integer not null,
	x smallint not null,
	y smallint not null,
	z smallint not null,
    primary key (pod_id, suffix_no),
    foreign key(pod_id) references pods(id),
    foreign key(object_id) references terrain_objects(id),
	foreign key(pod_id, x, y) referes pod_tiles(pod_id, x, y)
    );

create table pod_object_links
    (
    pod_id integer not null,
    suffix_no integer not null,
	destination_pod_id integer not null,
	destination_x smallint not null,
	destination_y smallint not null,
	destination_z smallint not null,
    primary key (pod_id, suffix_no),
	foreign key(pod_id, suffix_no) references pod_objects(pod_id, suffix_no)
    foreign key(destination_pod_id) references pods(id),
    );

