drop role if exists Reader;
drop role if exists Author;
drop role if exists Administrator;

create role Reader login;
create role Author login;
create role Administrator login;

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