create schema if not exists main;


drop table if exists main.note_collections;
drop table if exists interg_tests.teams_members;
drop table if exists interg_tests.teamss_sections;

drop table if exists main.texts;
drop table if exists main.images;
drop table if exists main.raw_datas;

drop table if exists main.notes;
drop table if exists main.collections;
drop table if exists main.sections;
drop table if exists interg_tests.teams;
drop table if exists main.users;

drop table if exists main.stat;


create table main.sections (
   id            serial       primary key,
   creation_date timestamptz  not null
);


create table interg_tests.teams (
    id                serial       primary key,
    name              varchar(255) not null unique,
    registration_date timestamptz  not null
);


create table main.users (
    id                serial       primary key,
    fio               varchar(255) not null,
    registration_date timestamptz  not null,
    login             varchar(255) not null unique,
    password          varchar(255) not null unique,
    role              int          default 0 check (role = 0 or role = 1 or role = 2)
);


create table main.notes (
    id                serial       primary key,
    access            int          not null check (access >= 0),
    name              varchar(255) not null unique,
    content_type      int          not null check (content_type = 1 OR content_type = 2),
    likes             int          default 0 check (likes >= 0),
    dislikes          int          default 0 check (dislikes >= 0),
    registration_date timestamptz  not null,
    owner_id          int          not null references main.users(id),
    section_id        int          not null references main.sections(id)
);


create table main.collections (
    id            serial       primary key,
    name          varchar(255) not null,
    creation_date timestamptz  not null,
    owner_id      int          not null references main.users(id)
);


create table main.note_collections (
    note_id       int not null references main.notes(id),
    collection_id int not null references main.collections(id),
    primary key (note_id, collection_id)
);


create table interg_tests.teams_members (
    team_id int not null references interg_tests.teams(id),
    user_id int not null references main.users(id),
    primary key (team_id, user_id)
);


create table interg_tests.teamss_sections (
    team_id    int not null unique references interg_tests.teams(id),
    section_id int not null unique references main.sections(id),
    primary key (team_id, section_id)
);


create table main.texts (
    id       serial     primary key,
    data     bytea      not null,
    note_id  int        not null references main.notes(id),
    file_ext varchar(25) default 'txt'
);


create table main.images (
    id       serial      primary key,
    data     bytea       not null,
    note_id  int         not null references main.notes(id),
    file_ext varchar(25) not null
);


create table main.raw_datas (
    id       serial      primary key,
    data     bytea       not null,
    note_id  int         not null references main.notes(id),
    file_ext varchar(25) not null
);


create table main.stat (
    id                         serial      primary key,
    total_users                int         default 0,
    total_readers              int         default 0,
    total_authors              int         default 0,
    total_admins               int         default 0,
    total_collections          int         default 0,
    total_teams                int         default 0,
    total_sections             int         default 0,
    total_notes                int         default 0,
    total_open_notes           int         default 0,
    total_close_notes          int         default 0,
    total_notes_in_collections int         default 0,
    stat_date                  timestamptz default now()
);


------------------------------------------------------------------------------------------------------------------------
-- Roles
------------------------------------------------------------------------------------------------------------------------


drop role if exists Reader;
drop role if exists Author;
drop role if exists Administrator;

create role Reader login;
create role Author login;
create role Administrator login;

grant select on main.notes, main.texts, main.images, main.raw_datas, interg_tests.teams to Reader;
grant select on interg_tests.teams, interg_tests.teams_members, main.sections, interg_tests.teamss_sections to Reader;
grant select, update, insert, delete on main.collections, main.note_collections to Reader;
grant usage, select on all sequences in schema main to Reader;

grant select, update, insert, delete on main.notes, main.texts, main.images, main.raw_datas to Author;
grant select on interg_tests.teams, interg_tests.teams_members, main.sections, interg_tests.teamss_sections to Author;
grant select, update, insert, delete on main.collections, main.note_collections to Author;
grant usage, select on all sequences in schema main to Author;

grant select, update, insert, delete on main.notes, main.texts, main.images, main.raw_datas to Administrator;
grant select, update, insert, delete on main.collections, main.note_collections to Administrator;
grant select, update, insert, delete on interg_tests.teamss_sections to Administrator;
grant select, update, insert, delete on interg_tests.teams to Administrator;
grant select, update, insert, delete on interg_tests.teams_members to Administrator;
grant select, update, insert, delete on main.sections to Administrator;
grant select, update, insert, delete on main.users to Administrator;
grant usage, select on all sequences in schema main to Administrator;


------------------------------------------------------------------------------------------------------------------------
-- Triggers
------------------------------------------------------------------------------------------------------------------------


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table users
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_user_trigger()
returns trigger as $$
declare
    cur_date_time timestamptz := now();
begin
    -- Копируем последнюю строку из таблицы main.stat
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, cur_date_time
    from main.stat
    order by stat_date desc
    limit 1;

    -- В зависимости от операции обновляем соответствующее поле
    if TG_OP = 'INSERT' then
        update main.stat
        set total_users = total_users + 1,
            total_readers = total_readers + case when new.role = 0 then 1 else 0 end,
            total_authors = total_authors + case when new.role = 1 then 1 else 0 end,
            total_admins = total_admins + case when new.role = 2 then 1 else 0 end
        where stat_date = cur_date_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_users = total_users - 1,
            total_readers = total_readers - case when old.role = 0 then 1 else 0 end,
            total_authors = total_authors - case when old.role = 1 then 1 else 0 end,
            total_admins = total_admins - case when old.role = 2 then 1 else 0 end
        where stat_date = cur_date_time;

        return old;

    elsif TG_OP = 'UPDATE' then
        -- Если изменяется роль пользователя
        if new.role != old.role then
            update main.stat
            set total_readers = total_readers + case when new.role = 0 then 1 else 0 end - case when old.role = 0 then 1 else 0 end,
                total_authors = total_authors + case when new.role = 1 then 1 else 0 end - case when old.role = 1 then 1 else 0 end,
                total_admins = total_admins + case when new.role = 2 then 1 else 0 end - case when old.role = 2 then 1 else 0 end
            where stat_date = cur_date_time;

            return new;
        end if;
    end if;

    return new;
end;
$$ language plpgsql;

drop trigger if exists stat_user_trigger on main.users;
create trigger stat_user_trigger
    after insert or delete or update on main.users
    for each row
execute function main.func_stat_user_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table notes
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_note_trigger()
returns trigger as $$
declare
    cur_date_time timestamptz := now();
begin
    -- Копируем последнюю строку из таблицы main.stat
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, cur_date_time
    from main.stat
    order by stat_date desc
    limit 1;

    -- В зависимости от операции обновляем соответствующее поле
    if TG_OP = 'INSERT' then
        update main.stat
        set total_notes = total_notes + 1,
            total_open_notes = total_open_notes + case when new.access = 1 then 1 else 0 end,
            total_close_notes = total_close_notes + case when new.access = 0 then 1 else 0 end
        where stat_date = cur_date_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_notes = total_notes - 1,
            total_open_notes = total_open_notes - case when old.access = 1 then 1 else 0 end,
            total_close_notes = total_close_notes - case when old.access = 0 then 1 else 0 end
        where stat_date = cur_date_time;

        return old;

    elsif TG_OP = 'UPDATE' then
        -- Если изменяется доступность заметки
        if new.access != old.access then
            update main.stat
            set total_open_notes = total_open_notes + case when new.access = 1 then 1 else 0 end - case when old.access = 1 then 1 else 0 end,
                total_close_notes = total_close_notes + case when new.access = 0 then 1 else 0 end - case when old.access = 0 then 1 else 0 end
            where stat_date = cur_date_time;

            return new;
        end if;
    end if;

    return new;
end;
$$ language plpgsql;

drop trigger if exists stat_note_trigger on main.notes;
create trigger stat_note_trigger
    after insert or delete or update on main.notes
    for each row
execute function main.func_stat_note_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table sections
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_sections_trigger()
returns trigger as $$
declare
    cur_date_time timestamptz := now();
begin
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, cur_date_time
    from main.stat
    order by stat_date desc
    limit 1;

    if TG_OP = 'INSERT' then
        update main.stat
        set total_sections = total_sections + 1
        where stat_date = cur_date_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_sections = total_sections - 1
        where stat_date = cur_date_time;

        return old;

    end if;

    return new;
end;
$$ language plpgsql;

drop trigger if exists stat_section_trigger on main.sections;
create trigger stat_sections_trigger
    after insert or update or delete on main.sections
    for each row
execute procedure main.func_stat_sections_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table collections
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_collections_trigger()
returns trigger as $$
declare
    cur_date_time timestamptz := now();
begin
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, cur_date_time
    from main.stat
    order by stat_date desc
    limit 1;

    if TG_OP = 'INSERT' then
        update main.stat
        set total_collections = total_collections + 1
        where stat_date = cur_date_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_collections = total_collections - 1
        where stat_date = cur_date_time;

        return old;

    end if;

    return new;
end;
$$ language plpgsql;

drop trigger if exists stat_collections_trigger on main.collections;
create trigger stat_collections_trigger
    after insert or update or delete on main.collections
    for each row
execute procedure main.func_stat_collections_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table data
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_teams_trigger()
returns trigger as $$
declare
    cur_date_time timestamptz := now();
begin
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, cur_date_time
    from main.stat
    order by stat_date desc
    limit 1;

    if TG_OP = 'INSERT' then
        update main.stat
        set total_teams = total_teams + 1
        where stat_date = cur_date_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_teams = total_teams - 1
        where stat_date = cur_date_time;

        return old;

    end if;

    return new;
end;
$$ language plpgsql;

drop trigger if exists stat_teams_trigger on interg_tests.teams;
create trigger stat_teams_trigger
    after insert or update or delete on interg_tests.teams
    for each row
execute procedure main.func_stat_teams_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table team_members
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_team_members_trigger()
returns trigger as $$
declare
    cur_date_time timestamptz := now();
begin
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, cur_date_time
    from main.stat
    order by stat_date desc
    limit 1;

    if TG_OP = 'INSERT' then
        update main.stat
        set total_teams = total_teams + 1
        where stat_date = cur_date_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_teams = total_teams - 1
        where stat_date = cur_date_time;

        return old;

    end if;

    return new;
end;
$$ language plpgsql;

drop trigger if exists stat_team_members_trigger on interg_tests.teams_members;
create trigger stat_team_members_trigger
    after insert or update or delete on interg_tests.teams_members
    for each row
execute procedure main.func_stat_team_members_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table teams_sections
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_teams_sections_trigger()
returns trigger as $$
declare
    cur_date_time timestamptz := now();
begin
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, cur_date_time
    from main.stat
    order by stat_date desc
    limit 1;

    if TG_OP = 'INSERT' then
        update main.stat
        set total_sections = total_sections + 1
        where stat_date = cur_date_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_sections = total_sections - 1
        where stat_date = cur_date_time;

        return old;

    end if;

    return new;
end;
$$ language plpgsql;

drop trigger if exists stat_teams_sections_trigger on interg_tests.teamss_sections;
create trigger stat_teams_sections_trigger
    after insert or update or delete on interg_tests.teamss_sections
    for each row
execute procedure main.func_stat_teams_sections_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table note_collections
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_note_collections_trigger()
returns trigger as $$
declare
    cur_date_time timestamptz := now();
begin
    -- Вставляем новую строку в таблицу stat, используя последнюю запись
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, cur_date_time
    from main.stat
    order by stat_date desc
    limit 1;

    if TG_OP = 'INSERT' then
        -- Увеличиваем счетчик total_notes_in_collections на 1
        update main.stat
        set total_notes_in_collections = total_notes_in_collections + 1
        where stat_date = cur_date_time;

        return new;

    elsif TG_OP = 'DELETE' then
        -- Уменьшаем счетчик total_notes_in_collections на 1
        update main.stat
        set total_notes_in_collections = total_notes_in_collections - 1
        where stat_date = cur_date_time;

        return old;

    end if;

    return new; -- Возвращаем new по умолчанию
end;
$$ language plpgsql;


create trigger stat_note_collections_trigger
    after insert or delete on main.note_collections
    for each row
execute procedure main.func_stat_note_collections_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Necessary inserts
------------------------------------------------------------------------------------------------------------------------


insert into main.users (fio, registration_date, login, password, role) values
    ('mainadmin', now(), 'mainadminlogin', 'mainadminpassword', 2);

insert into main.sections (creation_date) values
    (now());

insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date) values
    (1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, now());