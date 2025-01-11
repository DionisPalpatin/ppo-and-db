package data_access

// Users queries
const (
	getUserByIDQuery = `
	SELECT * 
	FROM %s.users
    WHERE id = $1;
	`
	getUserByStringQuery = `
	SELECT * 
	FROM %s.users
    WHERE login = $1 OR fio = $1;
	`
	getAllUsersQuery = `
	SELECT * 
	FROM %s.users;
	`
	addUserQuery = `
	INSERT INTO %s.users (fio, registration_date, login, password, role) VALUES
		($1, $2, $3, $4, $5)
	returning id;
	`
	deleteUserQuery = `
	call %s.delete_user($1);
	`
	updateUserQuery = `
	UPDATE %s.users SET fio = $1, registration_date = $2, login = $3, password = $4, role = $5
	WHERE id = $6;
	`
)

// Team queries
const (
	getTeamByIDQuery = `
	SELECT *
	FROM %s.teams
	WHERE id = $1;
	`
	getTeamByNameQuery = `
	SELECT *
	FROM %s.teams
	WHERE name = $1;
	`
	// getTeamBySectionIDQuery = `
	// SELECT *
	// FROM %s.teams t
	// JOIN %s.teams_sections ts ON t.id = ts.team_id
	// WHERE ts.section_id = $1;
	// `
	getAllTeamsQuery = `
	SELECT * FROM %s.teams;
	`
	addTeamQuery = `
	INSERT INTO %s.teams (name, registration_date) VALUES 
		($1, $2)
	returning id;
	`
	deleteTeamQuery = `
	call %s.delete_team($1);
	`
	addUserToTeamQuery = `
	INSERT INTO %s.teams_members (team_id, user_id) values
	($1, $2);
	`
	deleteUserFromTeamQuery = `
	DELETE FROM %s.teams_members
    WHERE user_id = $1 AND team_id = $2;
	`
	updateTeamQuery = `
	UPDATE %s.teams SET name = $1, registration_date = $2
	WHERE id = $3;
	`
	getTeamMembersQuery = `
	SELECT 
        u.id, 
        u.fio, 
        u.registration_date, 
        u.login, 
        u.password, 
        u.role
    FROM %s.users u
    JOIN %s.teams_members tm ON u.id = tm.user_id
    WHERE tm.team_id = $1;
	`
	getUserTeamQuery = `
	SELECT id, name, registration_date
	FROM %s.teams t
	JOIN %s.teams_members tm ON t.id = tm.team_id
	WHERE tm.user_id = $1;
	`
)

// Note queries
const (
	getNoteByIDQuery = `
	SELECT *
	FROM %s.notes
	WHERE id = $1;
	`
	getNoteByNameQuery = `
	SELECT *
	FROM %s.notes
	WHERE name = $1;
	`
	getNoteTextContent = `
	SELECT data, file_ext
	FROM %s.texts
	WHERE note_id = $1;
	`
	getNoteImageContent = `
	SELECT data, file_ext
	FROM %s.images
	WHERE note_id = $1;
	`
	getNoteRawDataContent = `
	SELECT data, file_ext
	FROM %s.raw_datas
	WHERE note_id = $1;
	`
	getAllNotesQuery = `
	SELECT *
	FROM %s.notes;
	`
	getAllPublicNotesQuery = `
	SELECT *
	FROM %s.notes
	WHERE access = 1;
	`
	addNoteInfoQuery = `
	INSERT INTO %s.notes (access, name, content_type, likes, dislikes, registration_date, owner_id, section_id) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id;
	`
	addNoteTextQuery = `
	INSERT INTO %s.texts (data, file_ext, note_id) VALUES
	($1, $2, $3);
	`
	addNoteImageQuery = `
	INSERT INTO %s.images (data, file_ext, note_id) VALUES
	($1, $2, $3);
	`
	addNoteRawDataQuery = `
	INSERT INTO %s.raw_datas (data, file_ext, note_id) VALUES
	($1, $2, $3);
	`
	deleteNoteQuery = `
	call %s.delete_note($1);
	`
	updateNoteTextContentQuery = `
	UPDATE %s.texts SET (data, file_ext) = ($1, $2)
	WHERE note_id = $3;
	`
	updateNoteImageContentQuery = `
	UPDATE %s.images SET (data, file_ext) = ($1, $2)
	WHERE note_id = $3;
	`
	updateNoteRawContentQuery = `
	UPDATE %s.raw_datas SET (data, file_ext) = ($1, $2)
	WHERE note_id = $3;
	`
	updateNoteInfoQuery = `
	UPDATE %s.notes SET
	    access = $1,
	    name = $2,
	    content_type = $3,
	    likes = $4,
	    dislikes = $5,
	    registration_date = $6,
		owner_id = $7,
	    section_id = $8
	WHERE id = $9;
	`
	addNoteToCollectionQuery = `
	INSERT INTO %s.note_collections (note_id, collection_id) VALUES
		($1, $2);
	`
	deleteNoteFromCollectionQuery = `
	DELETE FROM %s.note_collections
	WHERE note_id = $1 AND collection_id = $2;
	`
)

// Collection queries
const (
	getCollectionByIDQuery = `
	SELECT *
	FROM %s.collections
	WHERE id = $1;
`
	getCollectionByNameQuery = `
	SELECT *
	FROM %s.collections
	WHERE name = $1;
`
	getAllCollectionsQuery = `
	SELECT *
	FROM %s.collections;
`
	getAllUserCollectionsQuery = `
	SELECT c.*
	FROM %s.collections c
		JOIN %s.note_collections nc ON c.id = nc.collection_id
		JOIN %s.notes n ON nc.note_id = n.id
	WHERE n.owner_id = $1;
`
	addCollectionQuery = `
	INSERT INTO %s.collections (name, creation_date, owner_id) VALUES
		($1, $2, $3)
	returning id;
`
	deleteCollectionQuery1 = `
	DELETE FROM %s.note_collections
	WHERE collection_id = $1;
`
	deleteCollectionQuery2 = `
	DELETE FROM %s.collections
	WHERE id = $1;
`
	updateCollectionQuery = `
	UPDATE %s.collections SET name = $1, creation_date = $2
	WHERE id = $3;
`
	getAllNotesInCollectionQuery = `
	SELECT n.id, n.access, n.name, n.content_type, n.likes, n.dislikes, 
       n.registration_date, n.owner_id, n.section_id
	FROM %s.notes n
		JOIN %s.note_collections nc ON n.id = nc.note_id
	WHERE nc.collection_id = $1
`
)

// Section queries
const (
	getSectionByIDQuery = `
	SELECT *
	FROM %s.sections
	WHERE id = $1;
`
	getSectionByTeamNameQuery = `
	SELECT id, creation_date
	FROM %s.sections s
		JOIN %s.teams_sections ts ON s.id = ts.section_id
	WHERE t.name = $1;
`
	getAllSectionsQuery = `
	SELECT *
	FROM %s.sections;
`
	addSectionQuery = `
	INSERT INTO %s.sections (creation_date) VALUES ($1)
	returning id;
`
	addSectionToTeamQuery = `
	INSERT INTO %s.teams_sections (team_id, section_id) VALUES
		($2, $3);
`
	deleteSectionQuery = `
	DELETE FROM %s.sections
	WHERE id = $1;

	DELETE FROM %s.teams_sections
	WHERE section_id = $1;
`
	updateSectionQuery = `
	UPDATE %s.sections SET creation_date = $1
	WHERE id = $2;
`
	getAllNotesInSectionQuery = `
	SELECT *
	FROM %s.notes
	WHERE section_id = $1;
`
	addNoteToSectionQuery = `
	UPDATE %s.notes SET section_id = $1
	WHERE id = $2;
`
	deleteNoteFromSectionQuery = `
	UPDATE %s.notes SET section_id = 0
	WHERE id = $1 AND section_id = $2;
`
)
