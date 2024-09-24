CREATE TABLE IF NOT EXISTS main.users (
                                          id UInt64,
                                          fio String,
                                          registration_date DateTime,
                                          login String,
                                          password String,
                                          role Int32
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS main.teams (
                                          id UInt64,
                                          name String,
                                          registration_date DateTime
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS main.sections (
                                             id UInt64,
                                             creation_date DateTime
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS main.notes (
                                          id UInt64,
                                          access Int32,
                                          name String,
                                          content_type Int32,
                                          likes Int32 DEFAULT 0,
                                          dislikes Int32 DEFAULT 0,
                                          registration_date DateTime,
                                          owner_id UInt64,
                                          section_id UInt64
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS main.collections (
                                                id UInt64,
                                                name String,
                                                creation_date DateTime,
                                                owner_id UInt64
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS main.note_collections (
                                                     note_id UInt64,
                                                     collection_id UInt64
) ENGINE = MergeTree()
ORDER BY (note_id, collection_id);

CREATE TABLE IF NOT EXISTS main.team_members (
                                                 team_id UInt64,
                                                 user_id UInt64
) ENGINE = MergeTree()
ORDER BY (team_id, user_id);

CREATE TABLE IF NOT EXISTS main.teams_sections (
                                                   team_id UInt64,
                                                   section_id UInt64
) ENGINE = MergeTree()
ORDER BY (team_id, section_id);

CREATE TABLE IF NOT EXISTS main.texts (
                                          id UInt64,
                                          data String,
                                          note_id UInt64
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS main.images (
                                           id UInt64,
                                           data String,
                                           note_id UInt64
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS main.raw_datas (
                                              id UInt64,
                                              data String,
                                              note_id UInt64
) ENGINE = MergeTree()
ORDER BY id;

INSERT INTO main.users (id, fio, registration_date, login, password, role) VALUES
    (1, 'ivanov ivan', '2006-01-02 15:04:05', 'adminlogin', 'adminpassword', 2);
