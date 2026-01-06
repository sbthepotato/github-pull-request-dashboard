package db_pkg

import (
	"context"
	"database/sql"
	"regexp"

	_ "modernc.org/sqlite"
)

type TitleRegex struct {
	TitleRegexId   *int    `json:"title_regex_id,omitempty"`
	RegexPattern   *string `json:"regex_pattern,omitempty"`
	Link           *string `json:"link,omitempty"`
	RepositoryName *string `json:"repository_name,omitempty"`
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

	_, err = db.ExecContext(
		ctx,
		`create table if not exists title_regex_repository (
			title_regex_repository_id integer primary key not null,
			title_regex_id integer not null,
			repository_name text not null,
			foreign key (title_regex_id) references title_regex(title_regex_id),
			foreign key (repository_name) references repository(name)
		)`,
	)
	if err != nil {
		return err
	}

	return nil

}

func (titleRegex *TitleRegex) init() {
	titleRegex.TitleRegexId = new(int)
	titleRegex.RegexPattern = new(string)
	titleRegex.Link = new(string)
	titleRegex.RepositoryName = new(string)
}

/**** public ****/
func UpsertTitleRegex(ctx context.Context, db *sql.DB, titleRegexes []*TitleRegex) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query, err := tx.PrepareContext(
		ctx,
		`insert into title_regex (
			title_regex_id,
			regex_pattern,
			link,
			repository_name
		) values (
			?,
			?,
			?,
			?
		) on conflict (title_regex_id) do update set
			regex_pattern = excluded.regex_pattern,
			link = excluded.link,
			repository_name = excluded.repository_name`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	for _, titleRegex := range titleRegexes {
		if *titleRegex.RegexPattern != "" && *titleRegex.Link != "" {

			_, err := regexp.Compile(*titleRegex.RegexPattern)
			if err != nil {
				return err
			}

			_, err = query.ExecContext(
				ctx,
				titleRegex.TitleRegexId,
				titleRegex.RegexPattern,
				titleRegex.Link,
				titleRegex.RepositoryName)
			if err != nil {
				tx.Rollback()
				return err
			}
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
			link,
			repository_name
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
			titleRegex.RepositoryName,
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

func DeleteTitleRegex(ctx context.Context, db *sql.DB, titleRegexId int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.QueryContext(ctx, `delete from title_regex_repository where title_regex_id = ?`, titleRegexId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.QueryContext(ctx, `delete from title_regex where title_regex_id = ?`, titleRegexId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
