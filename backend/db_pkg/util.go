package db_pkg

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDatabase(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "github-pull-request-dashboard.db")
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

	return db, nil
}
