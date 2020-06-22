create table agency
(
	agency_id serial not null
		constraint agency_pk
			primary key,
	agency_name varchar not null,
	created_on timestamp not null,
	slots integer default 0
);

create unique index agency_agency_name_uindex
	on agency (agency_name);

create table person
(
	person_id serial not null
		constraint person_pk
			primary key,
	agency_id integer not null
		constraint person_agency_agency_id_fk
			references agency
				on update cascade on delete cascade,
	email varchar not null,
	password_hash varchar not null,
	first_name varchar not null,
	last_name varchar not null,
	created_on timestamp not null,
	deleted boolean default false not null
);

create unique index person_email_uindex
	on person (email);

create table mission
(
	mission_id serial not null
		constraint mission_pk
			primary key,
	person_id integer not null
		constraint mission_person_person_id_fk
			references person
				on update cascade on delete restrict,
	agency_id integer not null
		constraint mission_agency_agency_id_fk
			references agency
				on update cascade on delete restrict,
	mission_name varchar not null,
	mission_state integer default 1 not null,
	created_on timestamp not null,
	updated_on timestamp not null,
	scheduled_start_date timestamp,
	scheduled_end_date timestamp,
	bbox_ul geometry not null,
	bbox_lr geometry not null
);

create unique index mission_mission_name_uindex
	on mission (mission_name);

create table flight
(
	flight_id serial not null
		constraint flight_pk
			primary key,
	mission_id integer not null
		constraint flight_mission_mission_id_fk
			references mission
				on update cascade on delete cascade,
	flight_state integer default 1 not null,
	engines_on_time timestamp,
	takeoff_time timestamp,
	landing_time timestamp,
	summary text default ''::text not null,
	created_on timestamp not null
);

create table file
(
	file_id serial not null
		constraint file_pk
			primary key,
	person_id integer not null
		constraint file_person_person_id_fk
			references person
				on update cascade on delete restrict,
	file_name text,
	file_type text not null,
	file_size bigint default 0 not null,
	hash_256 bytea not null,
	created_on timestamp not null,
	updated_on timestamp not null
);

create table mission_file
(
	mission_file_id serial not null
		constraint mission_file_pk
			primary key,
	mission_id integer not null
		constraint mission_file_mission_mission_id_fk
			references mission
				on update cascade on delete cascade,
	file_id integer not null
		constraint mission_file_file_file_id_fk
			references file
				on update cascade on delete cascade,
	created_on timestamp not null
);

create table agency_file
(
	agency_file_id serial not null
		constraint agency_file_pk
			primary key,
	agency_id integer not null
		constraint agency_file_agency_agency_id_fk
			references agency
				on update cascade on delete cascade,
	file_id integer not null
		constraint agency_file_file_file_id_fk
			references file
				on update cascade on delete cascade,
	created_on timestamp not null
);


create table file_downloads
(
	file_stats_id serial not null
		constraint file_downloads_pk
			primary key,
	file_id integer not null
		constraint file_downloads_file_file_id_fk
			references file
				on update cascade on delete cascade,
	person_id integer not null
		constraint file_downloads_person_person_id_fk
			references person
				on update cascade on delete cascade,
	created_on timestamp not null
);

create table spatial_ref_sys
(
	srid integer not null
		constraint spatial_ref_sys_pkey
			primary key
		constraint spatial_ref_sys_srid_check
			check ((srid > 0) AND (srid <= 998999)),
	auth_name varchar(256),
	auth_srid integer,
	srtext varchar(2048),
	proj4text varchar(2048)
);

create table mission_event
(
	mission_event_id serial not null
		constraint mission_event_pk
			primary key,
	event_type integer not null,
	payload jsonb default '{}'::jsonb not null,
	created_on timestamp
);

