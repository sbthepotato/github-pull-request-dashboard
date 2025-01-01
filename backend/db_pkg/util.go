package db_pkg

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDatabase(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "github-pull-request-dashboard.sqlite3")
	if err != nil {
		return nil, err
	}

	err = initRepositoryTable(ctx, db)
	if err != nil {
		return nil, err
	}

	err = initTeamTable(ctx, db)
	if err != nil {
		return nil, err
	}

	err = initUserTable(ctx, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

/*
turn sql null int to regular int pointer
*/
func nullIntToPtr(value sql.NullInt64) *int {
	if value.Valid {
		result := int(value.Int64)
		return &result
	}
	return nil
}

/*
turn sql null string to regular string pointer
*/
func nullStringToPtr(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}
