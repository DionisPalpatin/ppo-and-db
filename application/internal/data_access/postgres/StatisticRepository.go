package postgres

import (
	"context"
	_ "database/sql"
	_ "errors"
	"fmt"

	_ "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
)

func (sr *StatisticRepository) CountMarks() int {
	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT SUM(likes) AS total_likes, SUM(dislikes) AS total_dislikes FROM %s.notes;", schemaName)
	ctx := context.Background()

	likes := 0
	dislikes := 0

	err := db.QueryRowContext(ctx, query).Scan(&likes, &dislikes)
	if err != nil {
		return -1
	}

	return likes + dislikes
}

func (sr *StatisticRepository) CountUsers() int {
	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.users", schemaName)
	ctx := context.Background()

	res := 0

	err := db.QueryRowContext(ctx, query).Scan(&res)
	if err != nil {
		return -1
	}

	return res
}

func (sr *StatisticRepository) CountTeams() int {
	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.teams", schemaName)
	ctx := context.Background()

	res := 0

	err := db.QueryRowContext(ctx, query).Scan(&res)
	if err != nil {
		return -1
	}

	return res
}

func (sr *StatisticRepository) CountSections() int {
	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.sections", schemaName)
	ctx := context.Background()

	res := 0

	err := db.QueryRowContext(ctx, query).Scan(&res)
	if err != nil {
		return -1
	}

	return res
}

func (sr *StatisticRepository) CountNotes() int {
	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.notes", schemaName)
	ctx := context.Background()

	res := 0

	err := db.QueryRowContext(ctx, query).Scan(&res)
	if err != nil {
		return -1
	}

	return res
}

func (sr *StatisticRepository) CountCollections() int {
	db := sr.DbConfigs.DB
	schemaName := sr.DbConfigs.SchemaName
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.collections", schemaName)
	ctx := context.Background()

	res := 0

	err := db.QueryRowContext(ctx, query).Scan(&res)
	if err != nil {
		return -1
	}

	return res
}
