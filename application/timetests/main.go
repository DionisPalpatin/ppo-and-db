package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// ctx := context.Background()
	//
	// connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	"localhost", 5432, "postgres", "postgrespassword", "NotebookApp")
	//
	// db, err := sql.Open("postgres", connectionString)
	// if err != nil {
	// 	log.Fatalf("Unable to connect to database: %v", err)
	// }
	// defer db.Close()

	select_res := [4][6]time.Duration{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	}
	select_res_ind := [4][6]time.Duration{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	}
	insert_res := [4][2]time.Duration{
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
	}
	delete_res := [4][2]time.Duration{
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
	}
	update_res := [4][2]time.Duration{
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
	}

	repeats := 1

	for i := 1; i < 5; i++ {
		ctx := context.Background()

		connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			"localhost", 5432, "postgres", "postgrespassword", "NotebookApp")

		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}

		err = refreshTables(ctx, db)
		if err != nil {
			print(err.Error())
			return
		}

		err = generateTestData(ctx, db, i)
		if err != nil {
			print(err.Error())
			return
		}

		select_res[i-1][0] = testJoinQueries(ctx, db, repeats)
		select_res[i-1][1] = testOrderByQueries(ctx, db, repeats)
		select_res[i-1][2] = testGroupByQueries(ctx, db, repeats)
		select_res[i-1][3] = testWhereQueries(ctx, db, repeats)
		select_res[i-1][4] = testGroupByOrderByQueries(ctx, db, repeats)
		select_res[i-1][5] = testSubquery(ctx, db, repeats)

		select_res_ind[i-1][0] = testJoinQueriesIndex(ctx, db, repeats)
		select_res_ind[i-1][1] = testOrderByQueriesIndex(ctx, db, repeats)
		select_res_ind[i-1][2] = testGroupByQueriesIndex(ctx, db, repeats)
		select_res_ind[i-1][3] = testWhereQueriesIndex(ctx, db, repeats)
		select_res_ind[i-1][4] = testGroupByOrderByQueriesIndex(ctx, db, repeats)
		select_res_ind[i-1][5] = testSubqueryIndex(ctx, db, repeats)

		delete_res[i-1][0] = testDeleteWithoutIndex(ctx, db, repeats)
		delete_res[i-1][1] = testDeleteWithIndex(ctx, db, repeats)

		update_res[i-1][0] = testUpdateWithoutIndex(ctx, db, repeats)
		update_res[i-1][1] = testUpdateWithIndex(ctx, db, repeats)

		insert_res[i-1][0] = testSingleInsert(ctx, db, 1000*i, repeats)
		insert_res[i-1][1] = testBigInsert(ctx, db, 1000*i, repeats)

		db.Close()
	}

	// Вывод значений
	fmt.Println("SELECT Results (for Python):")
	fmt.Println("select_res = [")
	for _, row := range select_res {
		fmt.Printf("    [%d, %d, %d, %d, %d, %d],\n", row[0].Microseconds(), row[1].Microseconds(), row[2].Microseconds(), row[3].Microseconds(), row[4].Microseconds(), row[5].Microseconds())
	}
	fmt.Println("]")

	fmt.Println("\nSELECT Results with Index (for Python):")
	fmt.Println("select_res_ind = [")
	for _, row := range select_res_ind {
		fmt.Printf("    [%d, %d, %d, %d, %d, %d],\n", row[0].Microseconds(), row[1].Microseconds(), row[2].Microseconds(), row[3].Microseconds(), row[4].Microseconds(), row[5].Microseconds())
	}
	fmt.Println("]")

	fmt.Println("\nINSERT Results (for Python):")
	fmt.Println("insert_res = [")
	for _, row := range insert_res {
		fmt.Printf("    [%d, %d],\n", row[0].Microseconds(), row[1].Microseconds())
	}
	fmt.Println("]")

	fmt.Println("\nDELETE Results (for Python):")
	fmt.Println("delete_res = [")
	for _, row := range delete_res {
		fmt.Printf("    [%d, %d],\n", row[0].Microseconds(), row[1].Microseconds())
	}
	fmt.Println("]")

	fmt.Println("\nUPDATE Results (for Python):")
	fmt.Println("update_res = [")
	for _, row := range update_res {
		fmt.Printf("    [%d, %d],\n", row[0].Microseconds(), row[1].Microseconds())
	}
	fmt.Println("]")

}

// ---------------------------------------------------------------------------------------------------------------------
// Generate functions
// ---------------------------------------------------------------------------------------------------------------------

func truncateTables(ctx context.Context, conn *sql.DB) error {
	queries := []string{
		"TRUNCATE counttime.note_collections CASCADE;",
		"TRUNCATE counttime.team_members CASCADE;",
		"TRUNCATE counttime.teams_sections CASCADE;",
		"TRUNCATE counttime.notes CASCADE;",
		"TRUNCATE counttime.collections CASCADE;",
		"TRUNCATE counttime.sections CASCADE;",
		"TRUNCATE counttime.teams CASCADE;",
		"TRUNCATE counttime.users CASCADE;",

		"TRUNCATE counttime.note_collections1 CASCADE;",
		"TRUNCATE counttime.team_members1 CASCADE;",
		"TRUNCATE counttime.teams_sections1 CASCADE;",
		"TRUNCATE counttime.notes1 CASCADE;",
		"TRUNCATE counttime.collections1 CASCADE;",
		"TRUNCATE counttime.sections1 CASCADE;",
		"TRUNCATE counttime.teams1 CASCADE;",
		"TRUNCATE counttime.users1 CASCADE;",
	}

	for _, query := range queries {
		_, err := conn.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	return nil
}

func deleteTables(ctx context.Context, conn *sql.DB) error {
	queries := []string{
		"create schema if not exists counttime;",
		"drop table if exists counttime.note_collections;",
		"drop table if exists counttime.team_members;",
		"drop table if exists counttime.teams_sections;",
		"drop table if exists counttime.notes;",
		"drop table if exists counttime.collections;",
		"drop table if exists counttime.sections;",
		"drop table if exists counttime.teams;",
		"drop table if exists counttime.users;",

		"drop table if exists counttime.note_collections1;",
		"drop table if exists counttime.team_members1;",
		"drop table if exists counttime.teams_sections1;",
		"drop table if exists counttime.notes1;",
		"drop table if exists counttime.collections1;",
		"drop table if exists counttime.sections1;",
		"drop table if exists counttime.teams1;",
		"drop table if exists counttime.users1;",
	}

	for _, query := range queries {
		_, err := conn.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	return nil
}

func createTables(ctx context.Context, conn *sql.DB) error {
	query := `
	create schema if not exists counttime;
	
	create table counttime.sections (
	   id            serial       primary key,
	   creation_date timestamptz  not null
	);
	create table counttime.teams (
	    id                serial       primary key,
	    name              varchar(255) not null unique,
	    registration_date timestamptz  not null
	);
	create table counttime.users (
	    id                serial       primary key,
	    fio               varchar(255) not null,
	    registration_date timestamptz  not null,
	    login             varchar(255) not null unique,
	    password          varchar(255) not null unique,
	    role              int          default 0 check (role = 0 or role = 1 or role = 2)
	);
	create table counttime.notes (
	    id                serial       primary key,
	    access            int          not null check (access >= 0),
	    name              varchar(255) not null unique,
	    content_type      int          not null check (content_type = 1 OR content_type = 2),
	    likes             int          default 0 check (likes >= 0),
	    dislikes          int          default 0 check (dislikes >= 0),
	    registration_date timestamptz  not null,
	    owner_id          int          not null references counttime.users(id),
	    section_id        int          not null references counttime.sections(id)
	);
	create table counttime.collections (
	    id            serial       primary key,
	    name          varchar(255) not null,
	    creation_date timestamptz  not null,
	    owner_id      int          not null references counttime.users(id)
	);
	create table counttime.note_collections (
	    note_id       int not null references counttime.notes(id),
	    collection_id int not null references counttime.collections(id),
	    primary key (note_id, collection_id)
	);
	create table counttime.team_members (
	    team_id int not null references counttime.teams(id),
	    user_id int not null references counttime.users(id),
	    primary key (team_id, user_id)
	);
	create table counttime.teams_sections (
	    team_id    int not null unique references counttime.teams(id),
	    section_id int not null unique references counttime.sections(id),
	    primary key (team_id, section_id)
	);



	
	create table counttime.sections1 (
	   id            serial       primary key,
	   creation_date timestamptz  not null
	);
	create table counttime.teams1 (
	    id                serial       primary key,
	    name              varchar(255) not null unique,
	    registration_date timestamptz  not null
	);
	create table counttime.users1 (
	    id                serial       primary key,
	    fio               varchar(255) not null,
	    registration_date timestamptz  not null,
	    login             varchar(255) not null unique,
	    password          varchar(255) not null unique,
	    role              int          default 0 check (role = 0 or role = 1 or role = 2)
	);
	create table counttime.notes1 (
	    id                serial       primary key,
	    access            int          not null check (access >= 0),
	    name              varchar(255) not null unique,
	    content_type      int          not null check (content_type = 1 OR content_type = 2),
	    likes             int          default 0 check (likes >= 0),
	    dislikes          int          default 0 check (dislikes >= 0),
	    registration_date timestamptz  not null,
	    owner_id          int          not null references counttime.users1(id),
	    section_id        int          not null references counttime.sections1(id)
	);
	create table counttime.collections1 (
	    id            serial       primary key,
	    name          varchar(255) not null,
	    creation_date timestamptz  not null,
	    owner_id      int          not null references counttime.users1(id)
	);
	create table counttime.note_collections1 (
	    note_id       int not null references counttime.notes1(id),
	    collection_id int not null references counttime.collections1(id),
	    primary key (note_id, collection_id)
	);
	create table counttime.team_members1 (
	    team_id int not null references counttime.teams1(id),
	    user_id int not null references counttime.users1(id),
	    primary key (team_id, user_id)
	);
	create table counttime.teams_sections1 (
	    team_id    int not null unique references counttime.teams1(id),
	    section_id int not null unique references counttime.sections1(id),
	    primary key (team_id, section_id)
	);
	`
	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to refresh tables: %w", err)
	}

	return nil
}

func refreshTables(ctx context.Context, conn *sql.DB) error {
	err := deleteTables(ctx, conn)
	if err != nil {
		return err
	}

	err = createTables(ctx, conn)
	if err != nil {
		return err
	}

	return nil
}

func generateTestData(ctx context.Context, conn *sql.DB, i int) error {
	// Генерация данных для всех таблиц

	numUsers := 1000 * i
	numTeams := 100 * i
	numSections := 100 * i
	numNotes := 10000 * i
	numCollections := 1000 * i

	// Генерация данных
	if err := generateUsers(ctx, conn, numUsers); err != nil {
		log.Fatalf("Failed to generate users: %v", err)
		return err
	}
	if err := generateTeams(ctx, conn, numTeams); err != nil {
		log.Fatalf("Failed to generate teams: %v", err)
		return err
	}
	if err := generateSections(ctx, conn, numSections); err != nil {
		log.Fatalf("Failed to generate sections: %v", err)
		return err
	}
	if err := generateNotes(ctx, conn, numNotes, numUsers, numSections); err != nil {
		log.Fatalf("Failed to generate notes: %v", err)
		return err
	}
	if err := generateCollections(ctx, conn, numCollections, numUsers); err != nil {
		log.Fatalf("Failed to generate collections: %v", err)
		return err
	}
	if err := generateNoteCollections(ctx, conn, numNotes, numCollections); err != nil {
		log.Fatalf("Failed to generate note_collections: %v", err)
		return err
	}
	if err := generateTeamMembers(ctx, conn, numTeams, numUsers); err != nil {
		log.Fatalf("Failed to generate team_members: %v", err)
		return err
	}
	if err := generateTeamsSections(ctx, conn, numTeams, numSections); err != nil {
		log.Fatalf("Failed to generate teams_sections: %v", err)
		return err
	}

	return nil
}

func generateUsers(ctx context.Context, conn *sql.DB, count int) error {
	var queryValues string

	for i := 0; i < count; i++ {
		fio := fmt.Sprintf("'User_%d'", i+1)
		login := fmt.Sprintf("'login%d'", i+1)
		password := fmt.Sprintf("'password%d'", i+1)
		role := rand.Intn(3)
		registrationDate := fmt.Sprintf("'%s'", time.Now().Format(time.RFC3339))

		// Добавляем значения в строку
		queryValues += fmt.Sprintf("(%s, %s, %s, %s, %d),", fio, registrationDate, login, password, role)
	}

	// Убираем последнюю запятую
	queryValues = queryValues[:len(queryValues)-1]

	query := fmt.Sprintf(`
		INSERT INTO counttime.users (fio, registration_date, login, password, role)
		VALUES %s`, queryValues)

	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert users: %w", err)
	}

	query = fmt.Sprintf(`
		INSERT INTO counttime.users1 (fio, registration_date, login, password, role)
		VALUES %s`, queryValues)

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert users: %w", err)
	}

	fmt.Printf("%d users generated.\n", count)
	return nil
}

func generateTeams(ctx context.Context, conn *sql.DB, count int) error {
	var queryValues string

	for i := 0; i < count; i++ {
		name := fmt.Sprintf("'Team_%d'", i+1)
		registrationDate := fmt.Sprintf("'%s'", time.Now().Format(time.RFC3339))

		// Добавляем значения в строку
		queryValues += fmt.Sprintf("(%s, %s),", name, registrationDate)
	}

	// Убираем последнюю запятую
	queryValues = queryValues[:len(queryValues)-1]

	query := fmt.Sprintf(`
		INSERT INTO counttime.teams (name, registration_date)
		VALUES %s`, queryValues)

	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert teams: %w", err)
	}

	query = fmt.Sprintf(`
		INSERT INTO counttime.teams1 (name, registration_date)
		VALUES %s`, queryValues)

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert teams: %w", err)
	}

	fmt.Printf("%d teams generated.\n", count)
	return nil
}

func generateSections(ctx context.Context, conn *sql.DB, count int) error {
	var queryValues string

	for i := 0; i < count; i++ {
		creationDate := fmt.Sprintf("'%s'", time.Now().Format(time.RFC3339))

		// Добавляем значения в строку
		queryValues += fmt.Sprintf("(%s),", creationDate)
	}

	// Убираем последнюю запятую
	queryValues = queryValues[:len(queryValues)-1]

	query := fmt.Sprintf(`
		INSERT INTO counttime.sections (creation_date)
		VALUES %s`, queryValues)

	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert sections: %w", err)
	}

	query = fmt.Sprintf(`
		INSERT INTO counttime.sections1 (creation_date)
		VALUES %s`, queryValues)

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert sections: %w", err)
	}

	fmt.Printf("%d sections generated.\n", count)
	return nil
}

func generateNotes(ctx context.Context, conn *sql.DB, count int, numUsers int, numSections int) error {
	var queryValues string

	for i := 0; i < count; i++ {
		name := fmt.Sprintf("'Note_%d'", i+1)
		contentType := rand.Intn(2) + 1
		access := rand.Intn(2)
		ownerID := rand.Intn(numUsers) + 1
		sectionID := rand.Intn(numSections) + 1
		likes := rand.Intn(100)
		dislikes := rand.Intn(50)
		registrationDate := fmt.Sprintf("'%s'", time.Now().Format(time.RFC3339))

		// Добавляем значения в строку
		queryValues += fmt.Sprintf("(%d, %s, %d, %d, %d, %s, %d, %d),", access, name, contentType, likes, dislikes, registrationDate, ownerID, sectionID)
	}

	// Убираем последнюю запятую
	queryValues = queryValues[:len(queryValues)-1]

	query := fmt.Sprintf(`
		INSERT INTO counttime.notes (access, name, content_type, likes, dislikes, registration_date, owner_id, section_id)
		VALUES %s`, queryValues)

	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert notes: %w", err)
	}

	query = fmt.Sprintf(`
		INSERT INTO counttime.notes1 (access, name, content_type, likes, dislikes, registration_date, owner_id, section_id)
		VALUES %s`, queryValues)

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert notes: %w", err)
	}


	fmt.Printf("%d notes generated.\n", count)
	return nil
}

func generateCollections(ctx context.Context, conn *sql.DB, count int, numUsers int) error {
	var queryValues string

	for i := 0; i < count; i++ {
		name := fmt.Sprintf("'Collection_%d'", i+1)
		ownerID := rand.Intn(numUsers) + 1
		creationDate := fmt.Sprintf("'%s'", time.Now().Format(time.RFC3339))

		// Добавляем значения в строку
		queryValues += fmt.Sprintf("(%s, %s, %d),", name, creationDate, ownerID)
	}

	// Убираем последнюю запятую
	queryValues = queryValues[:len(queryValues)-1]

	query := fmt.Sprintf(`
		INSERT INTO counttime.collections (name, creation_date, owner_id)
		VALUES %s`, queryValues)

	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert collections: %w", err)
	}

	query = fmt.Sprintf(`
		INSERT INTO counttime.collections1 (name, creation_date, owner_id)
		VALUES %s`, queryValues)

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to insert collections: %w", err)
	}

	fmt.Printf("%d collections generated.\n", count)
	return nil
}

func generateTeamMembers(ctx context.Context, conn *sql.DB, numTeams, numUsers int) error {
	var queryValues string

	// Будем использовать карту для отслеживания уникальных связей team_id -> user_id
	teamMemberMap := make(map[int]map[int]bool)

	for teamID := 1; teamID <= numTeams; teamID++ {
		if teamMemberMap[teamID] == nil {
			teamMemberMap[teamID] = make(map[int]bool)
		}

		// Определим случайное количество участников для каждой команды
		numMembers := numUsers / 10 // От 1 до 10% пользователей могут быть в команде

		for i := 0; i < numMembers; i++ {
			for {
				userID := rand.Intn(numUsers) + 1

				// Проверим, нет ли уже этой комбинации team_id и user_id
				if !teamMemberMap[teamID][userID] {
					teamMemberMap[teamID][userID] = true
					queryValues += fmt.Sprintf("(%d, %d),", teamID, userID)
					break // Выходим из внутреннего цикла после добавления
				}
			}
		}
	}

	// Убираем последнюю запятую
	if len(queryValues) > 0 {
		queryValues = queryValues[:len(queryValues)-1]
	}

	query := fmt.Sprintf(`
		INSERT INTO counttime.team_members (team_id, user_id)
		VALUES %s`, queryValues)

	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("Failed to insert into team_members: %v", err)
		return err
	}

	query = fmt.Sprintf(`
		INSERT INTO counttime.team_members1 (team_id, user_id)
		VALUES %s`, queryValues)

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("Failed to insert into team_members: %v", err)
		return err
	}

	fmt.Println("Generated team_members data")
	return nil
}

func generateNoteCollections(ctx context.Context, conn *sql.DB, numNotes, numCollections int) error {
	var queryValues string

	// Используем карту для отслеживания уникальных связей note_id -> collection_id
	noteCollectionMap := make(map[int]map[int]bool)

	for noteID := 1; noteID <= numNotes; noteID++ {
		if noteCollectionMap[noteID] == nil {
			noteCollectionMap[noteID] = make(map[int]bool)
		}

		// Определим случайное количество коллекций для каждой заметки
		numCollectionsForNote := rand.Intn(3) + 1 // От 1 до 3 коллекций

		for i := 0; i < numCollectionsForNote; i++ {
			for {
				collectionID := rand.Intn(numCollections) + 1

				// Проверим, нет ли уже этой комбинации note_id и collection_id
				if !noteCollectionMap[noteID][collectionID] {
					noteCollectionMap[noteID][collectionID] = true
					queryValues += fmt.Sprintf("(%d, %d),", noteID, collectionID)
					break // Выходим из внутреннего цикла после добавления
				}
			}
		}
	}

	// Убираем последнюю запятую
	if len(queryValues) > 0 {
		queryValues = queryValues[:len(queryValues)-1]
	}

	query := fmt.Sprintf(`
		INSERT INTO counttime.note_collections (note_id, collection_id)
		VALUES %s`, queryValues)

	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("Failed to insert into note_collections: %v", err)
		return err
	}

	query = fmt.Sprintf(`
		INSERT INTO counttime.note_collections1 (note_id, collection_id)
		VALUES %s`, queryValues)

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("Failed to insert into note_collections: %v", err)
		return err
	}

	fmt.Println("Generated note_collections data")
	return nil
}

func generateTeamsSections(ctx context.Context, conn *sql.DB, numTeams, numSections int) error {
	var queryValues []string

	// Используем карту для отслеживания занятых секций
	occupiedSections := make(map[int]bool)

	// Создаем транзакцию
	tx, err := conn.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback() // Откат в случае ошибки

	for teamID := 1; teamID <= numTeams; teamID++ {
		for {
			sectionID := rand.Intn(numSections) + 1

			// Проверим, свободен ли раздел
			if !occupiedSections[sectionID] {
				occupiedSections[sectionID] = true // Отметим раздел как занятый
				queryValues = append(queryValues, fmt.Sprintf("(%d, %d)", teamID, sectionID))
				break // Выходим из внутреннего цикла после добавления
			}

			// Если раздел занят, продолжаем искать
		}
	}

	// Если есть значения для вставки, выполняем запрос
	if len(queryValues) > 0 {
		query := fmt.Sprintf(`
			INSERT INTO counttime.teams_sections (team_id, section_id)
			VALUES %s`, strings.Join(queryValues, ","))

		_, err = tx.ExecContext(ctx, query)
		if err != nil {
			log.Fatalf("Failed to insert into teams_sections: %v", err)
			return err
		}

		query = fmt.Sprintf(`
			INSERT INTO counttime.teams_sections1 (team_id, section_id)
			VALUES %s`, strings.Join(queryValues, ","))

		_, err = tx.ExecContext(ctx, query)
		if err != nil {
			log.Fatalf("Failed to insert into teams_sections: %v", err)
			return err
		}
	}

	// Завершаем транзакцию
	if err = tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
		return err
	}

	fmt.Println("Generated teams_sections data")
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Test functions
// ---------------------------------------------------------------------------------------------------------------------

// --- Test inserts ready

func testSingleInsert(ctx context.Context, conn *sql.DB, count int, repeats int) time.Duration {
	var result_time time.Duration

	for j := 0; j < repeats; j++ {
		for i := 0; i < count; i++ {
			fio := fmt.Sprintf("User_%d", i+1)
			login := fmt.Sprintf("login%d", i+1)
			password := fmt.Sprintf("password%d", i+1)
			role := rand.Intn(3) // Роль 0 - Reader, 1 - Author, 2 - Admin

			start := time.Now()
			conn.ExecContext(ctx, `
				INSERT INTO counttime.users (fio, registration_date, login, password, role)
				VALUES ($1, $2, $3, $4, $5)`,
				fio, time.Now(), login, password, role)
			result_time += time.Since(start)
		}
	}
	result_time /= time.Duration(repeats)

	return result_time
}

func testBigInsert(ctx context.Context, conn *sql.DB, count int, repeats int) time.Duration {
	var result_time time.Duration

	for j := 0; j < repeats; j++ {
		query := `INSERT INTO counttime.users (fio, registration_date, login, password, role) VALUES `

		// Генерация данных для вставки
		for i := 0; i < count; i++ {
			fio := fmt.Sprintf("User_%d", i+1)
			login := fmt.Sprintf("login%d", i+1)
			password := fmt.Sprintf("password%d", i+1)
			role := rand.Intn(3) // Роль: 0 - Reader, 1 - Author, 2 - Admin

			// Добавляем текущую строку в запрос
			if i > 0 {
				query += ", " // Добавляем запятую для последующих строк
			}
			query += fmt.Sprintf("('%s', '%s', '%s', '%s', %d)", fio, time.Now().Format(time.RFC3339), login, password, role)
		}

		// Замеряем время выполнения одного большого запроса
		start := time.Now()
		conn.ExecContext(ctx, query)
		result_time += time.Since(start)
	}

	return result_time / time.Duration(repeats)
}

// --- Test deletes ready

func testDeleteWithoutIndex(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	for i := 0; i < repeats; i++ {
		start := time.Now()
		conn.ExecContext(ctx, `
			DELETE FROM counttime.users
			WHERE role = 1;
		`)
		res += time.Since(start)
	}

	return res / time.Duration(repeats)
}

func testDeleteWithIndex(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	conn.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS counttime.idx_my ON counttime.users1(role)`)

	for i := 0; i < repeats; i++ {
		start := time.Now()
		conn.ExecContext(ctx, `
			DELETE FROM counttime.users1
			WHERE role = 1;
		`)
		res += time.Since(start)
	}

	conn.ExecContext(ctx, `DROP INDEX IF EXISTS counttime.idx_my`)

	return res / time.Duration(repeats)
}

// --- Test updates ready

func testUpdateWithoutIndex(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	for i := 0; i < repeats; i++ {
		start := time.Now()
		conn.ExecContext(ctx, `
	        UPDATE counttime.users
	        SET fio = 'Updated User'
	        WHERE role = 2;
	        `)
		res += time.Since(start)
	}

	return res / time.Duration(repeats)
}

func testUpdateWithIndex(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	conn.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS counttime.idx_my ON counttime.users1(role)`)

	for i := 0; i < repeats; i++ {
		start := time.Now()
		conn.ExecContext(ctx, `
	        UPDATE counttime.users1
	        SET fio = 'Updated User'
	        WHERE role = 2;
	        `)
		res += time.Since(start)
	}

	conn.ExecContext(ctx, `DROP INDEX IF EXISTS counttime.idx_my`)

	return res / time.Duration(repeats)
}

// --- Test selects

func testJoinQueries(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	for i := 0; i < repeats; i++ {
		start := time.Now()

		rows, err := conn.QueryContext(ctx, `
			SELECT n.id, n.name, u.fio
			FROM counttime.notes n
			JOIN counttime.users u ON n.owner_id = u.id
			WHERE n.likes > 10;
		`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	return res / time.Duration(repeats)
}

func testOrderByQueries(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	for i := 0; i < repeats; i++ {
		start := time.Now()

		rows, err := conn.QueryContext(ctx, `
			SELECT id, name, likes
			FROM counttime.notes
			ORDER BY likes DESC LIMIT 100;
		`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	return res / time.Duration(repeats)
}

func testGroupByQueries(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	for i := 0; i < repeats; i++ {
		start := time.Now()

		// Группировка заметок по владельцам (пользователям) и подсчет количества заметок
		rows, err := conn.QueryContext(ctx, `
			SELECT owner_id, COUNT(*)
			FROM counttime.notes
			GROUP BY owner_id;
		`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	return res / time.Duration(repeats)
}

func testWhereQueries(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	for i := 0; i < repeats; i++ {
		start := time.Now()

		// Группировка заметок по владельцам (пользователям) и подсчет количества заметок
		rows, err := conn.QueryContext(ctx, `
			SELECT id, name, likes
			FROM counttime.notes
			WHERE likes > 50 AND access = 1;
		`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	return res / time.Duration(repeats)
}

func testGroupByOrderByQueries(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	for i := 0; i < repeats; i++ {
		start := time.Now()

		// Группировка заметок по владельцам (пользователям) и подсчет количества заметок
		rows, err := conn.QueryContext(ctx, `
			SELECT owner_id, COUNT(*) as note_count
			FROM counttime.notes
			GROUP BY owner_id
			ORDER BY note_count DESC;
		`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	return res / time.Duration(repeats)
}

func testSubquery(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	for i := 0; i < repeats; i++ {
		start := time.Now()

		// Группировка заметок по владельцам (пользователям) и подсчет количества заметок
		rows, err := conn.QueryContext(ctx, `
			SELECT fio, login
			FROM counttime.users
			WHERE id IN (
				SELECT owner_id
				FROM counttime.notes
				WHERE likes > 50
		)`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	return res / time.Duration(repeats)
}

// ---

func testJoinQueriesIndex(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	conn.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS counttime.idx_my ON counttime.notes1(owner_id);`)

	for i := 0; i < repeats; i++ {
		start := time.Now()

		rows, err := conn.QueryContext(ctx, `
			SELECT n.id, n.name, u.fio
			FROM counttime.notes1 n
			JOIN counttime.users1 u ON n.owner_id = u.id
			WHERE n.likes > 10;
		`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()
		res += time.Since(start)
	}

	conn.ExecContext(ctx, `DROP INDEX IF EXISTS counttime.idx_my`)

	return res / time.Duration(repeats)
}

func testOrderByQueriesIndex(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	conn.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS counttime.idx_my ON counttime.notes1(likes);`)

	for i := 0; i < repeats; i++ {
		start := time.Now()

		rows, err := conn.QueryContext(ctx, `
			SELECT id, name, likes
			FROM counttime.notes1
			ORDER BY likes DESC LIMIT 100;
		`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	conn.ExecContext(ctx, `DROP INDEX IF EXISTS counttime.idx_my`)

	return res / time.Duration(repeats)
}

func testGroupByQueriesIndex(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	conn.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS counttime.idx_my ON counttime.notes1(owner_id);`)

	for i := 0; i < repeats; i++ {
		start := time.Now()

		// Группировка заметок по владельцам (пользователям) и подсчет количества заметок
		rows, err := conn.QueryContext(ctx, `
			SELECT owner_id, COUNT(*)
			FROM counttime.notes1
			GROUP BY owner_id;
		`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	conn.ExecContext(ctx, `DROP INDEX IF EXISTS counttime.idx_my`)

	return res / time.Duration(repeats)
}

func testWhereQueriesIndex(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	conn.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS counttime.idx_my ON counttime.notes1(likes);`)

	for i := 0; i < repeats; i++ {
		start := time.Now()

		// Группировка заметок по владельцам (пользователям) и подсчет количества заметок
		rows, err := conn.QueryContext(ctx, `
			SELECT id, name, likes
			FROM counttime.notes1
			WHERE likes > 50 AND access = 1;
		`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	conn.ExecContext(ctx, `DROP INDEX IF EXISTS counttime.idx_my`)

	return res / time.Duration(repeats)
}

func testGroupByOrderByQueriesIndex(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	conn.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS counttime.idx_my ON counttime.notes1(owner_id);`)

	for i := 0; i < repeats; i++ {
		start := time.Now()

		// Группировка заметок по владельцам (пользователям) и подсчет количества заметок
		rows, err := conn.QueryContext(ctx, `
			SELECT owner_id, COUNT(*) as note_count
			FROM counttime.notes1
			GROUP BY owner_id
			ORDER BY note_count DESC;
		`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	conn.ExecContext(ctx, `DROP INDEX IF EXISTS counttime.idx_my`)

	return res / time.Duration(repeats)
}

func testSubqueryIndex(ctx context.Context, conn *sql.DB, repeats int) time.Duration {
	var res time.Duration

	conn.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS counttime.idx_my ON counttime.notes1(likes);`)

	for i := 0; i < repeats; i++ {
		start := time.Now()

		// Группировка заметок по владельцам (пользователям) и подсчет количества заметок
		rows, err := conn.QueryContext(ctx, `
			SELECT fio, login
			FROM counttime.users1
			WHERE id IN (
				SELECT owner_id
				FROM counttime.notes1
				WHERE likes > 50
		)`)

		if err != nil {
			log.Fatalf("Query failed: %s", err.Error())
			return 0
		}

		defer rows.Close()

		res += time.Since(start)
	}

	conn.ExecContext(ctx, `DROP INDEX IF EXISTS counttime.idx_my`)

	return res / time.Duration(repeats)
}
