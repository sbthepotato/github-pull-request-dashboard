package db_pkg

import (
	"context"
	"database/sql"

	"github.com/google/go-github/v72/github"
	_ "modernc.org/sqlite"
)

type Team struct {
	*github.Team
	RepositoryName *string `json:"repository_name,omitempty"`
	ReviewOrder    *int    `json:"review_order,omitempty"`
}

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
func (*Team) init() {
	team := new(Team)

	team.RepositoryName = new(string)
	team.ReviewOrder = new(int)

	team.Team = new(github.Team)
	team.Team.Slug = new(string)
	team.Team.Name = new(string)
	team.Team.HTMLURL = new(string)

}

/*
get all the team slugs
*/
func getTeamsSlugsAsSlice(ctx context.Context, db *sql.DB) ([]string, error) {

	result, err := db.QueryContext(
		ctx,
		`select team.slug
		from team
		order by team.slug asc`)

	if err != nil {
		return nil, err
	}
	defer result.Close()

	slugs := make([]string, 0)

	for result.Next() {

		slug := ""

		err := result.Scan(&slug)
		if err != nil {
			return nil, err
		}

		slugs = append(slugs, slug)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return slugs, nil
}

/**** public ****/

/*
create many team rows in single transaction
if a team exists in the db but is missing in the new slice of teams then it will be automatically deleted
*/
func CreateTeams(ctx context.Context, db *sql.DB, teams []*Team) error {
	var err error

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

	newSlugs := make([]string, 0)

	for _, team := range teams {

		newSlugs = append(newSlugs, *team.Slug)

		_, err := query.ExecContext(ctx, team.Slug, team.Name, team.HTMLURL)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	oldSlugs, err := getTeamsSlugsAsSlice(ctx, db)
	if err != nil {
		tx.Rollback()
		return err
	}

	deletedSlugs := findExtraElements(newSlugs, oldSlugs)

	for _, slug := range deletedSlugs {
		_, err = tx.QueryContext(ctx,
			`delete from user_team where team_slug = ?`,
			slug)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.QueryContext(ctx,
			`delete from team_review where team_slug = ?`,
			slug)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.QueryContext(ctx,
			`delete from team where team_slug = ?`,
			slug)
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
		_, err := query.ExecContext(ctx, team.Slug, team.RepositoryName, team.ReviewOrder, team.ReviewOrder)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.QueryContext(ctx,
		`delete from user_team where team_slug in (
		select team_slug from team_review where review_order = 0)`)

	_, err = tx.QueryContext(ctx,
		`delete from team_review where review_order = 0`)
	if err != nil {
		tx.Rollback()
		return err
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

		team := new(Team)
		team.init()

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

		team := new(Team)
		team.init()

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
