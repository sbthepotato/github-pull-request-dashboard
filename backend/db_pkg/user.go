package db_pkg

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/go-github/v72/github"
	_ "modernc.org/sqlite"
)

type User struct {
	*github.User
	Team *Team `json:"team,omitempty"`
}

/**** private ****/

/*
initialize the user and user_team table
*/
func initUserTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`create table if not exists user ( 
			login text primary key not null,
			name text,
			html_url text not null,
			avatar_url text not null
		)`,
	)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		`create table if not exists user_team ( 
			user_login text not null,
			repository_name text not null,
			team_slug text not null,
			primary key (user_login, repository_name),
			foreign key (user_login) references user(login),
			foreign key (repository_name) references repository(name),
			foreign key (team_slug) references team(slug)
		)`,
	)
	if err != nil {
		return err
	}

	return nil
}

/*
create empty user struct with db fields
*/
func (*User) init() {
	user := new(User)

	user.User = new(github.User)
	user.User.Login = new(string)
	user.User.Name = new(string)
	user.User.HTMLURL = new(string)
	user.User.AvatarURL = new(string)

	user.Team = new(Team)
	user.Team.ReviewOrder = new(int)

	user.Team.Team = new(github.Team)
	user.Team.Team.Slug = new(string)
	user.Team.Team.Name = new(string)

}

/*
get query used to search users
*/
func getUserQuery(ctx context.Context, db *sql.DB, repositoryName string) (*sql.Rows, error) {

	result, err := db.QueryContext(
		ctx,
		`select 
			user.login,
			user.name as user_name,
			user.html_url,
			user.avatar_url,
			team.slug,
			team.name as team_name,
			team_review.review_order
		from user
		left join user_team on user.login = user_team.user_login 
			and user_team.repository_name = ?
		left join team on user_team.team_slug = team.slug	
		left join team_review on user_team.team_slug = team_review.team_slug
			and team_review.repository_name = ?
		order by coalesce(user.name, user.login) asc`,
		repositoryName,
		repositoryName)

	return result, err
}

/*
get all the user logins
*/
func getUserLoginAsSlice(ctx context.Context, db *sql.DB) ([]string, error) {
	result, err := db.QueryContext(
		ctx,
		`select user.login
		from user
		order by user.login asc`)

	if err != nil {
		return nil, err
	}
	defer result.Close()

	logins := make([]string, 0)

	for result.Next() {

		login := ""

		err := result.Scan(&login)
		if err != nil {
			return nil, err
		}

		logins = append(logins, login)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return logins, nil
}

/**** public ****/

/*
create many user rows in single transaction
if a user exists in the db but is missing in the new slice the old user will automatically be deleted
*/
func CreateUsers(ctx context.Context, db *sql.DB, users []*User) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query, err := tx.PrepareContext(
		ctx,
		`insert or ignore into user (
			login,
			name, 
			html_url,
			avatar_url
			) values (
			?,
			?,
			?,
			?
			)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	newLogins := make([]string, 0)

	for _, user := range users {

		newLogins = append(newLogins, *user.User.Login)

		_, err := query.ExecContext(ctx, user.User.Login, user.User.Name, user.User.HTMLURL, user.User.AvatarURL)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	oldLogins, err := getUserLoginAsSlice(ctx, db)
	if err != nil {
		tx.Rollback()
		return err
	}

	deletedLogins := findExtraElements(newLogins, oldLogins)

	for _, login := range deletedLogins {
		_, err = tx.QueryContext(ctx,
			`delete from user_team where user_login = ?`,
			login)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.QueryContext(ctx,
			`delete from user where user_login = ?`,
			login)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

/*
upsert user team data into user_team table
*/
func UpsertUserTeams(ctx context.Context, db *sql.DB, userTeams map[string][]*github.User, repositoryName string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query, err := tx.PrepareContext(
		ctx,
		`insert into user_team(
			user_login,
			repository_name,
			team_slug
		) values (
		?,
		?,
		?) on conflict (user_login, repository_name) do update set
		team_slug = ?`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer query.Close()

	for teamSlug, users := range userTeams {
		for _, user := range users {
			_, err := query.ExecContext(ctx, user.Login, repositoryName, teamSlug, teamSlug)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

/*
get users in map where team is key
*/
func GetUsersAsTeamMap(ctx context.Context, db *sql.DB, repositoryName string) (map[string][]*User, error) {

	result, err := getUserQuery(ctx, db, repositoryName)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	users := make(map[string][]*User)

	for result.Next() {

		user := new(User)
		user.init()

		var (
			userName    sql.NullString
			teamName    sql.NullString
			teamSlug    sql.NullString
			reviewOrder sql.NullInt64
		)

		err := result.Scan(
			user.User.Login,
			&userName,
			user.User.HTMLURL,
			user.User.AvatarURL,
			&teamSlug,
			&teamName,
			&reviewOrder,
		)
		if err != nil {
			log.Println("are we here again")
			return nil, err
		}

		user.User.Name = nullStringToPtr(userName)
		user.Team.Slug = nullStringToPtr(teamSlug)
		user.Team.Name = nullStringToPtr(teamName)
		user.Team.ReviewOrder = nullIntToPtr(reviewOrder)

		mapTeamName := ""
		if user.Team.Name == nil {
			mapTeamName = "none"
		} else {
			mapTeamName = *user.Team.Name
		}

		users[mapTeamName] = append(users[mapTeamName], user)

	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

/*
get users in map where login is the map key
*/
func GetUsersAsLoginMap(ctx context.Context, db *sql.DB, repositoryName string) (map[string]*User, error) {

	result, err := getUserQuery(ctx, db, repositoryName)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	users := make(map[string]*User)

	for result.Next() {

		user := new(User)
		user.init()

		var (
			userName    sql.NullString
			teamName    sql.NullString
			teamSlug    sql.NullString
			reviewOrder sql.NullInt64
		)

		err := result.Scan(
			user.User.Login,
			&userName,
			user.User.HTMLURL,
			user.User.AvatarURL,
			&teamSlug,
			&teamName,
			&reviewOrder,
		)
		if err != nil {
			return nil, err
		}

		user.User.Name = nullStringToPtr(userName)
		user.Team.Slug = nullStringToPtr(teamSlug)
		user.Team.Name = nullStringToPtr(teamName)
		user.Team.ReviewOrder = nullIntToPtr(reviewOrder)

		users[*user.Login] = user

	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return users, nil

}
