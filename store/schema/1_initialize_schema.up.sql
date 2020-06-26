CREATE EXTENSION IF NOT EXISTS postgis;

create table agency
(
    agency_id   serial    not null
        constraint agency_pk
            primary key,
    agency_name varchar   not null,
    created_on  timestamp not null,
    slots       integer default 0,
    invite_key  varchar   not null
);

create unique index agency_agency_name_uindex
    on agency (agency_name);

create table person
(
    person_id     serial                not null
        constraint person_pk
            primary key,
    agency_id     integer               not null
        constraint person_agency_agency_id_fk
            references agency
            on update cascade on delete cascade,
    email         varchar               not null,
    password_hash varchar               not null,
    first_name    varchar               not null,
    last_name     varchar               not null,
    created_on    timestamp             not null,
    deleted       boolean default false not null,
    rpc_token varchar not null
);

create unique index person_email_uindex
    on person (email);

create table mission
(
    mission_id           serial            not null
        constraint mission_pk
            primary key,
    person_id            integer           not null
        constraint mission_person_person_id_fk
            references person
            on update cascade on delete restrict,
    agency_id            integer           not null
        constraint mission_agency_agency_id_fk
            references agency
            on update cascade on delete restrict,
    mission_name         varchar           not null,
    mission_state        integer default 1 not null,
    created_on           timestamp         not null,
    updated_on           timestamp         not null,
    scheduled_start_date timestamp,
    scheduled_end_date   timestamp,
    bbox_ul              geometry          not null,
    bbox_lr              geometry          not null
);

create unique index mission_mission_name_uindex
    on mission (mission_name);

create table flight
(
    flight_id       serial                   not null
        constraint flight_pk
            primary key,
    mission_id      integer                  not null
        constraint flight_mission_mission_id_fk
            references mission
            on update cascade on delete cascade,
    flight_state    integer default 1        not null,
    engines_on_time timestamp,
    takeoff_time    timestamp,
    landing_time    timestamp,
    summary         text    default ''::text not null,
    created_on      timestamp                not null
);

create table file
(
    file_id    serial           not null
        constraint file_pk
            primary key,
    person_id  integer          not null
        constraint file_person_person_id_fk
            references person
            on update cascade on delete restrict,
    file_name  text,
    file_type  text             not null,
    file_size  bigint default 0 not null,
    hash_256   bytea            not null,
    created_on timestamp        not null,
    updated_on timestamp        not null
);

create table mission_file
(
    mission_file_id serial    not null
        constraint mission_file_pk
            primary key,
    mission_id      integer   not null
        constraint mission_file_mission_mission_id_fk
            references mission
            on update cascade on delete cascade,
    file_id         integer   not null
        constraint mission_file_file_file_id_fk
            references file
            on update cascade on delete cascade,
    created_on      timestamp not null
);

create table agency_file
(
    agency_file_id serial    not null
        constraint agency_file_pk
            primary key,
    agency_id      integer   not null
        constraint agency_file_agency_agency_id_fk
            references agency
            on update cascade on delete cascade,
    file_id        integer   not null
        constraint agency_file_file_file_id_fk
            references file
            on update cascade on delete cascade,
    created_on     timestamp not null
);

create table file_downloads
(
    file_stats_id serial    not null
        constraint file_downloads_pk
            primary key,
    file_id       integer   not null
        constraint file_downloads_file_file_id_fk
            references file
            on update cascade on delete cascade,
    person_id     integer   not null
        constraint file_downloads_person_person_id_fk
            references person
            on update cascade on delete cascade,
    created_on    timestamp not null
);

create table mission_event
(
    mission_event_id serial                    not null
        constraint mission_event_pk
            primary key,
    event_type       integer                   not null,
    payload          jsonb default '{}'::jsonb not null,
    created_on       timestamp,
    mission_id       integer                   not null
        constraint mission_event_mission_mission_id_fk
            references mission
            on update cascade on delete cascade
);

create table message
(
    message_id      serial                                not null
        constraint table_name_pk
            primary key,
    author_id       integer default 0                     not null,
    author_name     varchar,
    message_subject varchar default ''::character varying not null,
    message_body    text                                  not null,
    created_on      timestamp                             not null
);

create table person_inbox
(
    inbox_id   serial               not null
        constraint person_inbox_pk
            primary key,
    person_id  integer              not null
        constraint person_inbox_person_person_id_fk
            references person
            on update cascade on delete cascade,
    message_id integer
        constraint person_inbox_message_message_id_fk
            references message
            on update cascade on delete cascade,
    unread     boolean default true not null
);

create unique index person_inbox_person_id_message_id_uindex
    on person_inbox (person_id, message_id);
