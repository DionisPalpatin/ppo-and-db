create schema if not exists interg_tests;


drop table if exists interg_tests.note_collections CASCADE;
drop table if exists interg_tests.teams_members CASCADE;
drop table if exists interg_tests.teams_sections CASCADE;

drop table if exists interg_tests.texts CASCADE;
drop table if exists interg_tests.images CASCADE;
drop table if exists interg_tests.raw_datas CASCADE;

drop table if exists interg_tests.notes CASCADE;
drop table if exists interg_tests.collections CASCADE;
drop table if exists interg_tests.sections CASCADE;
drop table if exists interg_tests.teams CASCADE;
drop table if exists interg_tests.users CASCADE;

drop table if exists interg_tests.stat CASCADE;


create table interg_tests.sections (
   id            serial       primary key,
   creation_date timestamptz  not null
);


create table interg_tests.teams (
    id                serial       primary key,
    name              varchar(255) not null unique,
    registration_date timestamptz  not null
);


create table interg_tests.users (
    id                serial       primary key,
    fio               varchar(255) not null,
    registration_date timestamptz  not null,
    login             varchar(255) not null unique,
    password          varchar(255) not null unique,
    role              int          default 0 check (role = 0 or role = 1 or role = 2)
);


create table interg_tests.notes (
    id                serial       primary key,
    access            int          not null check (access >= 0),
    name              varchar(255) not null unique,
    content_type      int          not null check (content_type = 1 OR content_type = 2),
    likes             int          default 0 check (likes >= 0),
    dislikes          int          default 0 check (dislikes >= 0),
    registration_date timestamptz  not null,
    owner_id          int          not null references interg_tests.users(id),
    section_id        int          not null references interg_tests.sections(id)
);


create table interg_tests.collections (
    id            serial       primary key,
    name          varchar(255) not null,
    creation_date timestamptz  not null,
    owner_id      int          not null references interg_tests.users(id)
);


create table interg_tests.note_collections (
    note_id       int not null references interg_tests.notes(id),
    collection_id int not null references interg_tests.collections(id),
    primary key (note_id, collection_id)
);


create table interg_tests.teams_members (
    team_id int not null references interg_tests.teams(id),
    user_id int not null references interg_tests.users(id),
    primary key (team_id, user_id)
);


create table interg_tests.teams_sections (
    team_id    int not null unique references interg_tests.teams(id),
    section_id int not null unique references interg_tests.sections(id),
    primary key (team_id, section_id)
);


create table interg_tests.texts (
    id       serial     primary key,
    data     bytea      not null,
    note_id  int        not null references interg_tests.notes(id),
    file_ext varchar(25) default 'txt'
);


create table interg_tests.images (
    id       serial      primary key,
    data     bytea       not null,
    note_id  int         not null references interg_tests.notes(id),
    file_ext varchar(25) not null
);


create table interg_tests.raw_datas (
    id       serial      primary key,
    data     bytea       not null,
    note_id  int         not null references interg_tests.notes(id),
    file_ext varchar(25) not null
);


------------------------------------------------------------------------------------------------------------------------
-- хранимая процедура для удаления пользователя и всех связанных данных
------------------------------------------------------------------------------------------------------------------------
create or replace procedure interg_tests.delete_user(deleted_user_id integer)
as $$
begin
	delete from interg_tests.note_collections where collection_id in (select id from interg_tests.collections where owner_id = deleted_user_id);
	delete from interg_tests.collections where owner_id = deleted_user_id;
	delete from interg_tests.teams_members where user_id = deleted_user_id;
	delete from interg_tests.texts where note_id in (select id from interg_tests.notes where owner_id = deleted_user_id);
	delete from interg_tests.images where note_id in (select id from interg_tests.notes where owner_id = deleted_user_id);
	delete from interg_tests.raw_datas where note_id in (select id from interg_tests.notes where owner_id = deleted_user_id);
	delete from interg_tests.notes where owner_id = deleted_user_id;
	delete from interg_tests.users where id = deleted_user_id;
end;
$$ language plpgsql;




create or replace procedure interg_tests.delete_team(deleted_team_id integer)
as $$
declare
	team_sec_id int;
begin
	-- select section_id into team_sec_id
	-- from interg_tests.teams_sections
	-- where team_id = deleted_team_id;

	-- update interg_tests.notes set section_id = 0
	-- where section_id = team_sec_id;

	-- delete from interg_tests.teams_sections where team_id = deleted_team_id;
	delete from interg_tests.teams_members where team_id = deleted_team_id;
	-- delete from interg_tests.sections where id = team_sec_id;
	delete from interg_tests.teams where id = deleted_team_id;
end;
$$ language plpgsql;


------------------------------------------------------------------------------------------------------------------------
-- хранимая процедура для удаления заметки и всех связанных данных
------------------------------------------------------------------------------------------------------------------------
create or replace procedure interg_tests.delete_note(deleted_note_id integer)
as $$
begin
  delete from interg_tests.texts where note_id = deleted_note_id;
  delete from interg_tests.images where note_id = deleted_note_id;
  delete from interg_tests.raw_datas where note_id = deleted_note_id;
  delete from interg_tests.note_collections where note_id = deleted_note_id;
  delete from interg_tests.notes where id = deleted_note_id;
end;
$$ language plpgsql;


------------------------------------------------------------------------------------------------------------------------
-- Necessary inserts
------------------------------------------------------------------------------------------------------------------------

insert into interg_tests.sections (creation_date)
values
    (now());


insert into interg_tests.users (fio, registration_date, login, password, role)
values
    ('user 1', now(), 'user1', 'password1', 2),
    ('user 2', now(), 'user2', 'password2', 2),
    ('user 3', now(), 'user3', 'password3', 2),
    ('user 4', now(), 'user4', 'password4', 1),
    ('user 5', now(), 'user5', 'password5', 1),
    ('user 6', now(), 'user6', 'password6', 1),
    ('user 7', now(), 'user7', 'password7', 1),
    ('user 8', now(), 'user8', 'password8', 0),
    ('user 9', now(), 'user9', 'password9', 0),
    ('user 10', now(), 'user10', 'password10', 0),
    ('user 11', now(), 'user11', 'password11', 0),
    ('user 12', now(), 'user12', 'password12', 0);


INSERT INTO interg_tests.teams (name, registration_date)
VALUES
    ('Alpha', '2023-04-13 00:00:00'),
    ('Beta', '2023-08-22 00:00:00'),
    ('Gamma', '2023-12-05 00:00:00'),
    ('Delta', '2023-04-25 00:00:00'),
    ('Epsilon', '2023-10-19 00:00:00'),
    ('Zeta', '2023-02-09 00:00:00'),
    ('Eta', '2023-01-23 00:00:00'),
    ('Theta', '2023-04-14 00:00:00'),
    ('Iota', '2023-10-29 00:00:00'),
    ('Kappa', '2023-07-22 00:00:00'),
    ('Lambda', '2023-08-04 00:00:00'),
    ('Mu', '2023-03-29 00:00:00');


INSERT INTO interg_tests.teams_members (user_id, team_id)
VALUES
    (1, 1),
    (2, 2),
    (3, 3),
    (5, 5);


INSERT INTO interg_tests.notes (access, name, content_type, registration_date, owner_id, section_id)
VALUES
    (1, 'Note1', 1, now(), 1, 1),
    (2, 'Note2', 2, now(), 1, 1),
    (0, 'Note3', 1, now(), 2, 1),
    (3, 'Note4', 2, now(), 3, 1),
    (1, 'Note5', 1, now(), 1, 1),
    (1, 'Note6', 2, now(), 2, 1),
    (2, 'Note7', 1, now(), 3, 1),
    (0, 'Note8', 2, now(), 1, 1),
    (3, 'Note9', 1, now(), 2, 1),
    (1, 'Note10', 2, now(), 3, 1),
    (2, 'Note11', 1, now(), 1, 1),
    (0, 'Note12', 2, now(), 2, 1);


INSERT INTO interg_tests.texts (data, note_id, file_ext)
VALUES
    ('Sample text data 1', 1, 'txt'),
    ('Sample text data 2', 3, 'txt'),
    ('Sample text data 3', 5, 'txt'),
    ('Sample text data 4', 7, 'txt'),
    ('Sample text data 5', 9, 'txt'),
    ('Sample text data 6', 11, 'txt');


INSERT INTO interg_tests.images (data, note_id, file_ext)
VALUES
    ('Sample image data 1', 2, 'png'),
    ('Sample image data 2', 4, 'jpg'),
    ('Sample image data 3', 6, 'png'),
    ('Sample image data 4', 8, 'jpg'),
    ('Sample image data 5', 10, 'png'),
    ('Sample image data 6', 12, 'jpg');


INSERT INTO interg_tests.raw_datas (data, note_id, file_ext)
VALUES
    ('Sample raw data 1', 1, 'bin'),
    ('Sample raw data 2', 2, 'dat'),
    ('Sample raw data 3', 3, 'bin'),
    ('Sample raw data 4', 4, 'dat'),
    ('Sample raw data 5', 5, 'bin'),
    ('Sample raw data 6', 6, 'dat');


INSERT INTO interg_tests.collections (name, creation_date, owner_id)
VALUES
    ('Collection 1', now(), 1),
    ('Collection 2', now(), 2),
    ('Collection 3', now(), 3),
    ('Collection 4', now(), 1),
    ('Collection 5', now(), 2),
    ('Collection 6', now(), 3),
    ('Collection 7', now(), 1),
    ('Collection 8', now(), 2),
    ('Collection 9', now(), 3),
    ('Collection 10', now(), 1),
    ('Collection 11', now(), 2),
    ('Collection 12', now(), 3);


INSERT INTO interg_tests.note_collections (note_id, collection_id)
VALUES
    (5, 5),
    (7, 7);