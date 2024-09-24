create schema if not exists main;


drop table if exists main.note_collections;
drop table if exists main.team_members;
drop table if exists main.teams_sections;

drop table if exists main.texts;
drop table if exists main.images;
drop table if exists main.raw_datas;

drop table if exists main.notes;
drop table if exists main.collections;
drop table if exists main.sections;
drop table if exists main.teams;
drop table if exists main.users;


create table main.sections (
   id            serial       primary key,
   creation_date varchar(255) not null
);


create table main.teams (
    id                serial       primary key,
    name              varchar(255) not null unique,
    registration_date varchar(255) not null
);


create table main.users (
    id                serial       primary key,
    fio               varchar(255) not null,
    registration_date varchar(255) not null,
    login             varchar(255) not null unique,
    password          varchar(255) not null unique,
    role              int          not null check (role >= 0)
);


create table main.notes (
    id                serial       primary key,
    access            int          not null check (access >= 0),
    name              varchar(255) not null,
    content_type      int          not null check (content_type >= 0),
    likes             int          default 0 check (likes >= 0),
    dislikes          int          default 0 check (dislikes >= 0),
    registration_date varchar(255) not null,
    owner_id          int          not null references main.users(id),
    section_id        int          not null references main.sections(id)
);


create table main.collections (
    id            serial       primary key,
    name          varchar(255) not null,
    creation_date varchar(255) not null,
    owner_id      int          not null references main.users(id)
);


create table main.note_collections (
    note_id       int not null references main.notes(id),
    collection_id int not null references main.collections(id),
    primary key (note_id, collection_id)
);


create table main.team_members (
    team_id int not null references main.teams(id),
    user_id int not null references main.users(id),
    primary key (team_id, user_id)
);


create table main.teams_sections (
    team_id    int not null unique references main.teams(id),
    section_id int not null unique references main.sections(id),
    primary key (team_id, section_id)
);


create table main.texts (
    id      serial primary key,
    data    bytea  not null,
    note_id int    not null references main.notes(id)
);


create table main.images (
    id      serial primary key,
    data    bytea  not null,
    note_id int    not null references main.notes(id)
);


create table main.raw_datas (
    id      serial primary key,
    data    bytea  not null,
    note_id int    not null references main.notes(id)
);


-- Триггер на автовыдачу прав при регистрастрации или изменении данных пользователя
CREATE OR REPLACE FUNCTION update_user_role()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.role = 0 THEN
        -- Роль читателя (reader)
        EXECUTE format('ALTER ROLE user_%s SET ROLE reader', NEW.id);
    ELSIF NEW.role = 1 THEN
        -- Роль автора (author)
        EXECUTE format('ALTER ROLE user_%s SET ROLE author', NEW.id);
    ELSIF NEW.role = 2 THEN
        -- Роль администратора (admin)
        EXECUTE format('ALTER ROLE user_%s SET ROLE admin', NEW.id);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS update_user_role_trigger ON main.users;
CREATE TRIGGER update_user_role_trigger
    AFTER INSERT OR UPDATE ON main.users
    FOR EACH ROW
EXECUTE PROCEDURE update_user_role();


-- ROLLLLLLLLLLES
create role NonMember login;
create role Reader login;
create role Author login;
create role Administrator login;

grant select, insert, update on main.users to NonMember;
grant usage, select on all sequences in schema main to NonMember;

grant select on main.notes, main.texts, main.images, main.raw_datas, main.teams to Reader;
grant select on main.teams, main.team_members, main.sections, main.teams_sections to Reader;
grant select, update, insert, delete on main.collections, main.note_collections to Reader;
grant usage, select on all sequences in schema main to Reader;

grant select, update, insert, delete on main.notes, main.texts, main.images, main.raw_datas to Author;
grant select on main.teams, main.team_members, main.sections, main.teams_sections to Author;
grant select, update, insert, delete on main.collections, main.note_collections to Author;
grant usage, select on all sequences in schema main to Author;

grant select, update, insert, delete on main.notes, main.texts, main.images, main.raw_datas to Administrator;
grant select, update, insert, delete on main.collections, main.note_collections to Administrator;
grant select, update, insert, delete on main.teams_sections to Administrator;
grant select, update, insert, delete on main.teams to Administrator;
grant select, update, insert, delete on main.team_members to Administrator;
grant select, update, insert, delete on main.sections to Administrator;
grant select, update, insert, delete on main.users to Administrator;
grant usage, select on all sequences in schema main to Administrator;


insert into main.users (fio, registration_date, login, password, role) values
    ('ivanov ivan', '2006-01-02 15:04:05 -0700', 'adminlogin', 'adminpassword', 2);


insert into main.sections (creation_date) values
    ('0000-00-00 00:00:00 -0000');