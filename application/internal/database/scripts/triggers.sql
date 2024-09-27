-- Триггер на автовыдачу прав при регистрастрации или изменении данных пользователя
CREATE OR REPLACE FUNCTION main.update_user_role()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.role = 0 THEN
        -- Роль читателя (reader)
        EXECUTE format('ALTER ROLE user_%s SET ROLE Reader', NEW.id);
    ELSIF NEW.role = 1 THEN
        -- Роль автора (author)
        EXECUTE format('ALTER ROLE user_%s SET ROLE Author', NEW.id);
    ELSIF NEW.role = 2 THEN
        -- Роль администратора (admin)
        EXECUTE format('ALTER ROLE user_%s SET ROLE Admin', NEW.id);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;



DROP TRIGGER IF EXISTS update_user_role_trigger ON main.users;
CREATE TRIGGER update_user_role_trigger
    AFTER INSERT OR UPDATE ON main.users
    FOR EACH ROW
    EXECUTE PROCEDURE main.update_user_role();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table users
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_user_trigger()
    returns trigger as $$
declare
    latest_stat_date timestamptz;
    current_time timestamptz := now();
begin
    -- Получаем максимальную дату из таблицы stat
    select max(stat_date) into latest_stat_date from main.stat;

    -- Копируем последнюю строку из таблицы main.stat
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, current_time
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
        where stat_date = latest_stat_date;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_users = total_users - 1,
            total_readers = total_readers - case when old.role = 0 then 1 else 0 end,
            total_authors = total_authors - case when old.role = 1 then 1 else 0 end,
            total_admins = total_admins - case when old.role = 2 then 1 else 0 end
        where stat_date = latest_stat_date;

        return old;

    elsif TG_OP = 'UPDATE' then
        -- Если изменяется роль пользователя
        if new.role != old.role then
            update main.stat
            set total_readers = total_readers + case when new.role = 0 then 1 else 0 end - case when old.role = 0 then 1 else 0 end,
                total_authors = total_authors + case when new.role = 1 then 1 else 0 end - case when old.role = 1 then 1 else 0 end,
                total_admins = total_admins + case when new.role = 2 then 1 else 0 end - case when old.role = 2 then 1 else 0 end
            where stat_date = latest_stat_date;

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
    latest_stat_date timestamptz;
    current_time timestamptz := now();
begin
    -- Получаем максимальную дату из таблицы stat
    select max(stat_date) into latest_stat_date from main.stat;

    -- Копируем последнюю строку из таблицы main.stat
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, current_time
    from main.stat
    order by stat_date desc
    limit 1;

    -- В зависимости от операции обновляем соответствующее поле
    if TG_OP = 'INSERT' then
        update main.stat
        set total_notes = total_notes + 1,
            total_open_notes = total_open_notes + case when new.access = 1 then 1 else 0 end,
            total_close_notes = total_close_notes + case when new.access = 0 then 1 else 0 end
        where stat_date = latest_stat_date;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_notes = total_notes - 1,
            total_open_notes = total_open_notes - case when old.access = 1 then 1 else 0 end,
            total_close_notes = total_close_notes - case when old.access = 0 then 1 else 0 end
        where stat_date = latest_stat_date;

        return old;

    elsif TG_OP = 'UPDATE' then
        -- Если изменяется доступность заметки
        if new.access != old.access then
            update main.stat
            set total_open_notes = total_open_notes + case when new.access = 1 then 1 else 0 end - case when old.access = 1 then 1 else 0 end,
                total_close_notes = total_close_notes + case when new.access = 0 then 1 else 0 end - case when old.access = 0 then 1 else 0 end
            where stat_date = latest_stat_date;

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
    latest_stat_date timestamptz;
    current_time timestamptz := now();
begin
    select max(stat_date) into latest_stat_date from main.stat;

    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, current_time
    from main.stat
    where stat_date = latest_stat_date
    limit 1;

    if TG_OP = 'INSERT' then
        update main.stat
        set total_sections = total_sections + 1
        where stat_date = current_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_sections = total_sections - 1
        where stat_date = current_time;

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
    latest_stat_date timestamptz;
    current_time timestamptz := now();
begin
    select max(stat_date) into latest_stat_date from main.stat;

    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, current_time
    from main.stat
    where stat_date = latest_stat_date
    limit 1;

    if TG_OP = 'INSERT' then
        update main.stat
        set total_collections = total_collections + 1
        where stat_date = current_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_collections = total_collections - 1
        where stat_date = current_time;

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
-- Trigger on table teams
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_teams_trigger()
    returns trigger as $$
declare
    latest_stat_date timestamptz;
    current_time timestamptz := now();
begin
    select max(stat_date) into latest_stat_date from main.stat;

    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, current_time
    from main.stat
    where stat_date = latest_stat_date
    limit 1;

    if TG_OP = 'INSERT' then
        update main.stat
        set total_teams = total_teams + 1
        where stat_date = current_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_teams = total_teams - 1
        where stat_date = current_time;

        return old;

    end if;

    return new;
end;
$$ language plpgsql;

drop trigger if exists stat_teams_trigger on main.teams;
create trigger stat_teams_trigger
    after insert or update or delete on main.teams
    for each row
execute procedure main.func_stat_teams_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table team_members
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_team_members_trigger()
    returns trigger as $$
declare
    latest_stat_date timestamptz;
    current_time timestamptz := now();
begin
    select max(stat_date) into latest_stat_date from main.stat;

    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, current_time
    from main.stat
    where stat_date = latest_stat_date
    limit 1;

    if TG_OP = 'INSERT' then
        update main.stat
        set total_teams = total_teams + 1
        where stat_date = current_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_teams = total_teams - 1
        where stat_date = current_time;

        return old;

    end if;

    return new;
end;
$$ language plpgsql;

drop trigger if exists stat_team_members_trigger on main.team_members;
create trigger stat_team_members_trigger
    after insert or update or delete on main.team_members
    for each row
execute procedure main.func_stat_team_members_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table teams_sections
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_teams_sections_trigger()
    returns trigger as $$
declare
    latest_stat_date timestamptz;
    current_time timestamptz := now();
begin
    select max(stat_date) into latest_stat_date from main.stat;

    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, current_time
    from main.stat
    where stat_date = latest_stat_date
    limit 1;

    if TG_OP = 'INSERT' then
        update main.stat
        set total_sections = total_sections + 1
        where stat_date = current_time;

        return new;

    elsif TG_OP = 'DELETE' then
        update main.stat
        set total_sections = total_sections - 1
        where stat_date = current_time;

        return old;

    end if;

    return new;
end;
$$ language plpgsql;

drop trigger if exists stat_teams_sections_trigger on main.teams_sections;
create trigger stat_teams_sections_trigger
    after insert or update or delete on main.teams_sections
    for each row
execute procedure main.func_stat_teams_sections_trigger();


------------------------------------------------------------------------------------------------------------------------
-- Trigger on table note_collections
------------------------------------------------------------------------------------------------------------------------
create or replace function main.func_stat_note_collections_trigger()
    returns trigger as $$
declare
    latest_stat_date timestamptz;
    current_time timestamptz := now();
begin
    select max(stat_date) into latest_stat_date from main.stat;

    -- Вставляем новую строку в таблицу stat, используя последнюю запись
    insert into main.stat (total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, stat_date)
    select total_users, total_readers, total_authors, total_admins, total_collections, total_teams, total_sections, total_notes, total_open_notes, total_close_notes, total_notes_in_collections, current_time
    from main.stat
    where stat_date = latest_stat_date
    limit 1;

    if TG_OP = 'INSERT' then
        -- Увеличиваем счетчик total_notes_in_collections на 1
        update main.stat
        set total_notes_in_collections = total_notes_in_collections + 1
        where stat_date = current_time;

        return new;

    elsif TG_OP = 'DELETE' then
        -- Уменьшаем счетчик total_notes_in_collections на 1
        update main.stat
        set total_notes_in_collections = total_notes_in_collections - 1
        where stat_date = current_time;

        return old;

    end if;

    return new; -- Возвращаем new по умолчанию
end;
$$ language plpgsql;

create trigger stat_note_collections_trigger
    after insert or delete on main.note_collections
    for each row
execute procedure main.func_stat_note_collections_trigger();
