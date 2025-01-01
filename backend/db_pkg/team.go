package db_pkg

import (
	"context"
	"database/sql"

	"github.com/google/go-github/v68/github"
	_ "modernc.org/sqlite"
)

/**** private ****/

/*
initialize the team and team_review table
*/
func initTeamTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`create table if not exists team ( 
			slug text primary key not null,
			name text not null,
			html_url text not null
		)`,
	)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		`create table if not exists team_review ( 
			team_slug text not null,
			repository_name text not null,
			review_order integer not null,
			primary key (team_slug, repository_name),
			foreign key (team_slug) references team(slug),
			foreign key (repository_name) references repository(name)
		)`,
	)
	if err != nil {
		return err
	}

	return nil
}

/*
create empty team struct with db fields
*/
func initTeamStruct() *Team {
	team := new(Team)

	team.RepositoryName = new(string)
	team.ReviewOrder = new(int)

	team.Team = new(github.Team)
	team.Team.Slug = new(string)
	team.Team.Name = new(string)
	team.Team.HTMLURL = new(string)

	return team
}

/**** public ****/

/*
create many team rows in single transaction
*/
func CreateTeams(ctx context.Context, db *sql.DB, teams []*Team) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query, err := tx.PrepareContext(
		ctx,
		`insert or ignore into team (
			slug,
			name, 
			html_url
			) values (
			?,
			?,
			?
			)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	for _, team := range teams {

		_, err := query.ExecContext(ctx, team.Slug, team.Name, team.HTMLURL)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func UpsertTeamReviews(ctx context.Context, db *sql.DB, teams []*Team) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query, err := tx.PrepareContext(
		ctx,
		`insert into team_review(
			team_slug,
			repository_name,
			review_order
		) values (
		?,
		?,
		?) on conflict (team_slug, repository_name) do update set
		review_order = ?`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	for _, team := range teams {
		if *team.ReviewOrder > 0 {
			_, err := query.ExecContext(ctx, team.Slug, team.RepositoryName, team.ReviewOrder, team.ReviewOrder)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

/*
Get a list of teams with their reviews
*/
func GetTeams(ctx context.Context, db *sql.DB, repositoryName string) ([]*Team, error) {

	result, err := db.QueryContext(
		ctx,
		`select
			team.slug,
			team.name,
			team.html_url,
			coalesce(team_review.repository_name, ?),
			coalesce(team_review.review_order, 0)
		from team
		left join team_review on team.slug = team_review.team_slug 
			and team_review.repository_name = ?
		order by name asc`,
		repositoryName,
		repositoryName)

	if err != nil {
		return nil, err
	}
	defer result.Close()

	teams := make([]*Team, 0)

	for result.Next() {

		team := initTeamStruct()

		err := result.Scan(
			team.Slug,
			team.Name,
			team.HTMLURL,
			team.RepositoryName,
			team.ReviewOrder,
		)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return teams, nil

}

/*
get teams as map where team slug is key
*/
func GetTeamsAsMap(ctx context.Context, db *sql.DB, repositoryName string) (map[string]*Team, error) {

	result, err := db.QueryContext(
		ctx,
		`select
			team.slug,
			team.name,
			team.html_url,
			team_review.repository_name,
			team_review.review_order
		from team
		inner join team_review on team.slug = team_review.team_slug 
			and team_review.repository_name = ?
		order by name asc`,
		repositoryName,
		repositoryName)

	if err != nil {
		return nil, err
	}
	defer result.Close()

	teams := make(map[string]*Team)

	for result.Next() {

		team := initTeamStruct()

		err := result.Scan(
			team.Slug,
			team.Name,
			team.HTMLURL,
			team.RepositoryName,
			team.ReviewOrder,
		)
		if err != nil {
			return nil, err
		}

		teams[*team.Slug] = team
	}

	return teams, nil
}
