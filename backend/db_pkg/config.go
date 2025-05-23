package db_pkg

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

/**** private ****/

/*
initialize the config table
*/
func initConfigTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`create table if not exists title_regex (
			table_regex_id primary key not null,
			regex_pattern text not null,
			link text not null
		)`,
	)
	if err != nil {
		return err
	}

	return nil

}

/**** public ****/
func UpsertTitleRegex(ctx context.Context, db *sql.DB, titleRegexes []*TitleRegex) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query, err := tx.PrepareContext(
		ctx,
		`insert or update into title_regex (
		)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	for _, titleRegex := range titleRegexes {
		_, err := query.ExecContext(ctx, titleRegex.Table_regex_id, titleRegex.Regex_pattern, titleRegex.Link)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
