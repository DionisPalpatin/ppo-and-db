------------------------------------------------------------------------------------------------------------------------
-- хранимая процедура для удаления пользователя и всех связанных данных
------------------------------------------------------------------------------------------------------------------------
create or replace procedure main.delete_user(deleted_user_id integer)
as $$
begin
	delete from main.notes_collections where collection_id in (select id from main.collections where owner_id = deleted_user_id);
	delete from main.collections where owner_id = deleted_user_id;
	delete from main.team_members where user_id = deleted_user_id;
	delete from main.texts where note_id in (select id from main.notes where owner_id = deleted_user_id);
	delete from main.images where note_id in (select id from main.notes where owner_id = deleted_user_id);
	delete from main.raw_datas where note_id in (select id from main.notes where owner_id = deleted_user_id);
	delete from main.notes where owner_id = deleted_user_id;
	delete from main.users where id = deleted_user_id;
end;
$$ language plpgsql;




create or replace procedure main.delete_team(deleted_team_id integer)
as $$
declare
	team_sec_id int;
begin
	select section_id into team_sec_id
	from main.teams_sections
	where team_id = deleted_team_id;

	update main.notes set section_id = 0
	where section_id = team_sec_id;

	delete from main.teams_sections where team_id = deleted_team_id;
	delete from main.team_members where team_id = deleted_team_id;
	delete from main.sections where id = team_sec_id;
	delete from main.teams where id = deleted_team_id;
end;
$$ language plpgsql;


------------------------------------------------------------------------------------------------------------------------
-- хранимая процедура для удаления заметки и всех связанных данных
------------------------------------------------------------------------------------------------------------------------
create or replace procedure main.delete_note(deleted_note_id integer)
as $$
begin
  delete from main.texts where note_id = deleted_note_id;
  delete from main.images where note_id = deleted_note_id;
  delete from main.raw_datas where note_id = deleted_note_id;
  delete from main.notes_collections where note_id = deleted_note_id;
  delete from main.notes where id = deleted_note_id;
end;
$$ language plpgsql;