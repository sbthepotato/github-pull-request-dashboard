package db_pkg

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

type TitleRegex struct {
	TitleRegexId *int    `json:"title_regex_id,omitempty"`
	RegexPattern *string `json:"regex_pattern,omitempty"`
	Link         *string `json:"link,omitempty"`
}

/**** private ****/

/*
initialize the config table
*/
func initConfigTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`create table if not exists title_regex (
			title_regex_id integer primary key not null,
			regex_pattern text not null,
			link text not null
		)`,
	)
	if err != nil {
		return err
	}

	return nil

}

func (*TitleRegex) init() {
	TitleRegex := new(TitleRegex)

	TitleRegex.TitleRegexId = new(int)
	TitleRegex.RegexPattern = new(string)
	TitleRegex.Link = new(string)
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
		_, err := query.ExecContext(ctx, titleRegex.TitleRegexId, titleRegex.RegexPattern, titleRegex.Link)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetTitleRegexList(ctx context.Context, db *sql.DB) ([]*TitleRegex, error) {
	query, err := db.QueryContext(
		ctx,
		`select
			title_regex_id,
			regex_pattern,
			link
		from title_regex
		order by title_regex_id asc`,
	)

	if err != nil {
		return nil, err
	}

	result := make([]*TitleRegex, 0)

	for query.Next() {
		titleRegex := new(TitleRegex)
		titleRegex.init()

		err := query.Scan(
			titleRegex.TitleRegexId,
			titleRegex.RegexPattern,
			titleRegex.Link,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, titleRegex)

	}

	if err := query.Err(); err != nil {
		return nil, err
	}

	return result, nil

}
