-- for version 0.2.14

drop table pods;

create table pods
	(
	id int,
	pod_title text,
	editable boolean,
	primary key(id)
	);

insert into pods (id, pod_title, editable) values (1,'Login Pod',false);
insert into pods (id, pod_title, editable) values (2,'Three Tier Pod',false);
insert into pods (id, pod_title, editable) values (3,'Large Pod',false);
insert into pods (id, pod_title, editable) values (4,'Test Pod Four',true);
insert into pods (id, pod_title, editable) values (5,'Test Pod Five',true);
insert into pods (id, pod_title, editable) values (6,'Seed Pod',false);

