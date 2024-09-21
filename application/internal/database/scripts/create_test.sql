create schema if not exists test;


drop table if exists test.note_collections;
drop table if exists test.team_members;
drop table if exists test.teams_sections;

drop table if exists test.texts;
drop table if exists test.images;
drop table if exists test.raw_datas;

drop table if exists test.notes;
drop table if exists test.collections;
drop table if exists test.sections;
drop table if exists test.teams;
drop table if exists test.users;


create table test.sections (
                               id            serial       primary key,
                               creation_date varchar(255) not null
);


create table test.teams (
                            id                serial       primary key,
                            name              varchar(255) not null unique,
                            registration_date varchar(255) not null
);


create table test.users (
                            id                serial       primary key,
                            fio               varchar(255) not null,
                            registration_date varchar(255) not null,
                            login             varchar(255) not null unique,
                            password          varchar(255) not null unique,
                            role              int          not null check (role >= 0)
);


create table test.notes (
                            id                serial       primary key,
                            access            int          not null check (access >= 0),
                            name              varchar(255) not null,
                            content_type      int          not null check (content_type >= 0),
                            likes             int          default 0 check (likes >= 0),
                            dislikes          int          default 0 check (dislikes >= 0),
                            registration_date varchar(255) not null,
                            owner_id          int          not null references test.users(id),
                            section_id        int          not null references test.sections(id)
);


create table test.collections (
                                  id            serial       primary key,
                                  name          varchar(255) not null,
                                  creation_date varchar(255) not null,
                                  owner_id      int          not null references test.users(id)
);


create table test.note_collections (
                                       note_id       int not null references test.notes(id),
                                       collection_id int not null references test.collections(id),
                                       primary key (note_id, collection_id)
);


create table test.team_members (
    team_id int not null references test.teams(id),
    user_id int not null references test.users(id),
    primary key (team_id, user_id)
);


create table test.teams_sections (
    team_id    int not null unique references test.teams(id),
    section_id int not null unique references test.sections(id),
    primary key (team_id, section_id)
);


create table test.texts (
    id      serial primary key,
    data    bytea  not null,
    note_id int    not null references test.notes(id)
);


create table test.images (
     id      serial primary key,
     data    bytea  not null,
     note_id int    not null references test.notes(id)
);


create table test.raw_datas (
    id      serial primary key,
    data    bytea  not null,
    note_id int    not null references test.notes(id)
);


insert into test.users (fio, registration_date, login, password, role) values
   ('ivanov ivan', '05 Dec 2000', 'ivanovlogin', 'ivanovpassword', 0),
   ('petrov petr', '01 Dec 2003', 'petrovlogin', 'petrovpassword', 1),
   ('kostov kola', '02 Dec 1995', 'l', 'p', 2);