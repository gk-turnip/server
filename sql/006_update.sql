
-- for version 0.2.10

-- holds user preferences

create table pref_users
	(
	id int, -- user id
	screen_width smallint,
	screen_height smallint,
	zoom numeric(4,2), -- 0.5 to 2.0
	fps_display boolean, -- display frames per second
	pos_display boolean, -- display position
	pod_display boolean, -- display current pod name
	chat_display_lines smallint,
	background_music_volume_level numeric(5,2), -- 0% to 100%
	sound_effects_volume_level numeric(5,2), -- 0% to 100%
	primary key(id),
	foreign key(id) references users(id)
	);

