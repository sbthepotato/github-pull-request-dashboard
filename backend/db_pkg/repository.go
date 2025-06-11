package db_pkg

import (
	"context"
	"database/sql"

	"github.com/google/go-github/v72/github"
	_ "modernc.org/sqlite"
)

type Repository struct {
	*github.Repository
	Enabled *bool `json:"enabled,omitempty"`
}

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
			enabled integer not null
		)`,
	)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		`create index if not exists repository_enabled on repository(enabled)`)
	if err != nil {
		return err
	}

	return nil
}

/*
create a repository struct with the db fields
*/
func (*Repository) init() {
	repository := new(Repository)
	repository.Repository = new(github.Repository)
	repository.Name = new(string)
	repository.DefaultBranch = new(string)
	repository.HTMLURL = new(string)
	repository.Enabled = new(bool)

}

/**** public ****/

/*
create many repository rows in single transaction
*/
func CreateRepositories(ctx context.Context, db *sql.DB, repositories []*Repository) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query, err := tx.PrepareContext(
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
			?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	for _, repository := range repositories {

		_, err := query.ExecContext(ctx, repository.Name, repository.DefaultBranch, repository.HTMLURL, repository.Enabled)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

/*
set many repositories in a single database transaction
*/
func SetRepositories(ctx context.Context, db *sql.DB, repositories []*Repository) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query, err := tx.PrepareContext(
		ctx,
		`update repository set
			enabled = coalesce(?, enabled)
			where name = ?`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	for _, repository := range repositories {

		_, err := query.ExecContext(ctx, *repository.Enabled, *repository.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

/*
Get a list of repositories
*/
func GetRepositories(ctx context.Context, db *sql.DB, activeOnly bool) ([]*Repository, error) {

	result, err := db.QueryContext(
		ctx,
		`select 
			name,
			default_branch,
			html_url,
			enabled
		from repository
		where enabled = ?
			or ? = 0
		order by name asc`,
		activeOnly,
		activeOnly)

	if err != nil {
		return nil, err
	}
	defer result.Close()

	repositories := make([]*Repository, 0)

	for result.Next() {

		repository := new(Repository)
		repository.init()

		err := result.Scan(
			repository.Name,
			repository.DefaultBranch,
			repository.HTMLURL,
			repository.Enabled,
		)
		if err != nil {
			return nil, err
		}

		repositories = append(repositories, repository)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return repositories, nil

}
