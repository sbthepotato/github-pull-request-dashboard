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

	initRepositoryTable(ctx, db)

	return db, nil
}
