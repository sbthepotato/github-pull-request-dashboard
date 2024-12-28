package db_pkg

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDatabase(ctx context.Context) error {
	db, err := sql.Open("sqlite", "github-pull-request-dashboard.db")
	if err != nil {
		return err
	}

	initRepositoryTable(ctx, db)

	return nil
}
