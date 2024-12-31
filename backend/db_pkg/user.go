package db_pkg

import (
	"context"
	"database/sql"

	"github.com/google/go-github/v68/github"
	_ "modernc.org/sqlite"
)

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
			avatar_url text not null,
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
			team_name text not null,
			primary key (user_login, repository_name, team_name),
			foreign key (user_login) references user(login),
			foreign key (repository_name) references repository(name),
			foreign key (team_name) references team(name)
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
func initUserStruct() *User {
	user := new(User)

	user.User = new(github.User)
	user.User.Login = new(string)
	user.User.Name = new(string)
	user.User.HTMLURL = new(string)
	user.User.AvatarURL = new(string)

	user.Team = new(Team)
	user.Team.ReviewOrder = new(int)

	user.Team.Team = new(github.Team)
	user.Team.Team.Name = new(string)

	return user
}

/*
get query used to search users
*/
func getUserQuery(ctx context.Context, db *sql.DB, repositoryName string) (*sql.Rows, error) {

	result, err := db.QueryContext(
		ctx,
		`select 
			user.login,
			user.name,
			user.html_url,
			user.avatar_url,
			team_review.team_name,
			team_review.review_order
		from user
		left join user_team on user.login = user_team.user_login 
			and user_team.repository_name = ?
		left join team_review on user_team.team_name = team_review.team_name
			and team_review.repository_name = ?
		order by name asc`,
		repositoryName,
		repositoryName)

	return result, err
}

/**** public ****/

/*
create many user rows in single transaction
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

	for _, user := range users {

		_, err := query.ExecContext(ctx, user.User.Login, user.User.Name, user.User.HTMLURL, user.User.AvatarURL)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

/*
get users in map where team is key
*/
func GetUsersAsTeamMap(ctx context.Context, db *sql.DB, repositoryName string) (map[string]*User, error) {

	result, err := getUserQuery(ctx, db, repositoryName)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	users := make(map[string]*User, 0)

	for result.Next() {

		user := initUserStruct()

		// TODO this isnt done
		err := result.Scan(
			user.Name,
		)
		if err != nil {
			return nil, err
		}

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

	users := make(map[string]*User, 0)

	for result.Next() {

		user := initUserStruct()

		err := result.Scan(
			user.User.Login,
			user.User.Name,
			user.User.HTMLURL,
			user.User.AvatarURL,
			user.Team.Name,
			user.Team.ReviewOrder,
		)
		if err != nil {
			return nil, err
		}

		users[*user.Login] = user

	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return users, nil

}
