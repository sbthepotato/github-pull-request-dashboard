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
			name text primary key not null,
			html_url text not null
		)`,
	)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		`create table if not exists team_review ( 
			team_name text not null,
			repository_name text not null,
			review_order integer default 0,
			primary key (team_name, repository_name),
			foreign key (team_name) references team(name),
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
	team.Team = new(github.Team)
	team.Name = new(string)
	team.HTMLURL = new(string)
	team.RepositoryName = new(string)
	team.ReviewOrder = new(int)

	return team
}

/**** public ****/

/*
create a team row
*/
func CreateTeam(ctx context.Context, db *sql.DB, team *Team) error {

	_, err := db.ExecContext(
		ctx,
		`insert or ignore into team (
			name, 
			html_url
			) values (
			?,
			?
			)`,
		team.Name,
		team.HTMLURL,
	)
	if err != nil {
		return err
	}

	if team.RepositoryName != nil && *team.ReviewOrder != 0 {
		_, err := db.ExecContext(
			ctx,
			`insert or ignore into team_review (
				team_name,
				repository_name,
				review_order 
				) values (
				?,
				?,
				?
				)`,
			team.Name,
			team.RepositoryName,
			team.ReviewOrder,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

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
			name, 
			html_url
			) values (
			?,
			?
			)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	for _, team := range teams {

		_, err := query.ExecContext(ctx, team.Name, team.HTMLURL)
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
			team_name,
			repository_name,
			review_order
		) values (
		?,
		?,
		?) on conflict (team_name, repository_name) do update set
		review_order = ?`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	for _, team := range teams {
		if *team.ReviewOrder > 0 {
			_, err := query.ExecContext(ctx, team.Name, team.RepositoryName, team.ReviewOrder, team.ReviewOrder)
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
			team.name,
			team.html_url,
			coalesce(team_review.repository_name, ?),
			coalesce(team_review.review_order, 0)
		from team
		left join team_review on team.name = team_review.team_name 
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
