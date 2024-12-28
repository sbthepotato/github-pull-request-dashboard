package db_pkg

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

/**** private ****/

/*
initialize the repository table
*/
func initRepositoryTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`create table if not exists repository ( 
			name text primary key not null,
			default_branch text not null,
			html_url text not null,
			enabled integer
		)`,
	)
	if err != nil {
		return err
	}
	return nil
}

/**** public ****/

/*
create a repository
*/
func CreateRepository(ctx context.Context, repository *Repository) error {
	_, err := db.ExecContext(
		ctx,
		`insert or ignore into repository (
			name, 
			default_branch, 
			html_url, 
			enabled 
			) values (
			?,
			?,
			?,
			?)`,
		repository.Name,
		repository.DefaultBranch,
		repository.HTMLURL,
		repository.Enabled,
	)

	if err != nil {
		return err
	}

	return nil
}

func SetRepository(ctx context.Context, repository *Repository) error {
	_, err := db.ExecContext(
		ctx,
		`update repository set 
			enabled = ?
			where name = ?`,
		repository.Enabled,
		repository.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

/*
Get a list of repositories
*/
func GetRepositories(ctx context.Context, activeOnly bool) ([]*Repository, error) {

	result, err := db.QueryContext(
		ctx,
		`select 
		name,
		default_branch,
		html_url,
		enabled
		from repository
		where enabled = ?`,
		activeOnly)

	if err != nil {
		return nil, err
	}

	repositories := make([]*Repository, 0)

	for result.Next() {
		var repository Repository

		if err := result.Scan(&repository.Name,
			&repository.DefaultBranch,
			&repository.HTMLURL,
			&repository.Enabled); err != nil {

			return nil, err
		}
		repositories = append(repositories, &repository)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return repositories, nil

}
