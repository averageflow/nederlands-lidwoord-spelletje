insert into main.sqlite_master (type, name, tbl_name, rootpage, sql) values ('table', 'woord', 'woord', 2, 'CREATE TABLE woord
(
	id INTEGER not null
		constraint woord_pk
			primary key autoincrement,
	uid VARCHAR(36) not null,
	content VARCHAR(255) not null,
	created_at BIGINT not null
)');
insert into main.sqlite_master (type, name, tbl_name, rootpage, sql) values ('table', 'sqlite_sequence', 'sqlite_sequence', 3, 'CREATE TABLE sqlite_sequence(name,seq)');
insert into main.sqlite_master (type, name, tbl_name, rootpage, sql) values ('index', 'woord_content_uindex', 'woord', 4, 'CREATE UNIQUE INDEX woord_content_uindex
	on woord (content)');
insert into main.sqlite_master (type, name, tbl_name, rootpage, sql) values ('index', 'woord_id_uindex', 'woord', 5, 'CREATE UNIQUE INDEX woord_id_uindex
	on woord (id)');
insert into main.sqlite_master (type, name, tbl_name, rootpage, sql) values ('index', 'woord_uid_uindex', 'woord', 6, 'CREATE UNIQUE INDEX woord_uid_uindex
	on woord (uid)');
insert into main.sqlite_master (type, name, tbl_name, rootpage, sql) values ('table', 'lidwoord', 'lidwoord', 7, 'CREATE TABLE lidwoord
(
	id INTEGER not null
		constraint lidwoord_pk
			primary key autoincrement,
	uid VARCHAR(36) not null,
	content VARCHAR(255) not null,
	created_at BIGINT not null
)');
insert into main.sqlite_master (type, name, tbl_name, rootpage, sql) values ('table', 'woord_plural', 'woord_plural', 8, 'CREATE TABLE woord_plural
(
	id integer not null
		constraint woord_plural_pk
			primary key autoincrement,
	singular_id int not null
		constraint woord_plural_woord_id_fk
			references woord
				on update cascade on delete cascade,
	plural_id int not null
		constraint woord_plural_woord_id_fk_2
			references woord
				on update cascade on delete cascade
)');
insert into main.sqlite_master (type, name, tbl_name, rootpage, sql) values ('index', 'woord_plural_id_uindex', 'woord_plural', 9, 'CREATE UNIQUE INDEX woord_plural_id_uindex
	on woord_plural (id)');
insert into main.sqlite_master (type, name, tbl_name, rootpage, sql) values ('table', 'woord_lidwoord', 'woord_lidwoord', 10, 'CREATE TABLE woord_lidwoord
(
	id integer not null
		constraint woord_lidwoord_pk
			primary key autoincrement,
	woord_id int not null
		constraint woord_lidwoord_woord_id_fk
			references woord,
	lidwoord_id int not null
		constraint woord_lidwoord_lidwoord_id_fk
			references lidwoord
				on update cascade on delete cascade
)');
insert into main.sqlite_master (type, name, tbl_name, rootpage, sql) values ('index', 'woord_lidwoord_id_uindex', 'woord_lidwoord', 11, 'CREATE UNIQUE INDEX woord_lidwoord_id_uindex
	on woord_lidwoord (id)');
