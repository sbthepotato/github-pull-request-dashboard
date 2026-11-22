package db_pkg

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

type RepositoryGroup struct {
	RepositoryGroupId *int    `json:"repositoryGroupId,omitempty"`
	Name              *string `json:"name,omitempty"`
}

type RepositoryGroupRepository struct {
	RepositoryGroupId *int    `json:"repositoryGroupId,omitempty"`
	RepositoryName    *string `json:"repositoryName,omitempty"`
}

type RepositoryGroupRepositoryList struct {
	RepositoryGroup
	Repositories []RepositoryGroupRepository
}

/**** private ****/

func initRepositoryGroupTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`create table if not exists repository_group (
			repository_group_id integer primary key not null,
			name text not null
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		`create table if not exists repository_group_repository (
			repository_group_id integer not null,
			repository_name text not null,
			primary key (repository_group_id, repository_name),
			foreign key (repository_group_id) references repository_group(repository_group_id),
			foreign key (repository_name) references repository(name)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func (repositoryGroup *RepositoryGroup) init() {
	repositoryGroup.RepositoryGroupId = new(int)
	repositoryGroup.Name = new(string)
}

func (repositoryGroupRepository *RepositoryGroupRepository) init() {
	repositoryGroupRepository.RepositoryGroupId = new(int)
	repositoryGroupRepository.RepositoryName = new(string)
}

/**** public ****/
func UpsertRepositoryGroupRepositoryList(ctx context.Context, db *sql.DB, repositoryGroupRepositoryList []*RepositoryGroupRepositoryList) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	repository_group_query, err := tx.PrepareContext(
		ctx,
		`insert into repository_group (
			repository_group_id, 
			name
		) values (
			?, 
			?
		) on conflict(repository_group_id) do update set
			name = excluded.name
			returning repository_group_id;`,
	)
	if err != nil {
		return err
	}

	repository_group_repository_query, err := tx.PrepareContext(
		ctx,
		`insert into repository_group_repository (
			repository_group_id,
			repository_name
		) values (
			?,
			?
		) on conflict(repository_group_id, repository_name) do nothing;`,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, repositoryGroup := range repositoryGroupRepositoryList {
		_, err := repository_group_query.ExecContext(
			ctx,
			repositoryGroup.RepositoryGroupId,
			repositoryGroup.Name,
		)
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, repositoryGroupRepository := range repositoryGroup.Repositories {
			_, err := repository_group_repository_query.ExecContext(
				ctx,
				repositoryGroupRepository.RepositoryGroupId,
				repositoryGroupRepository.RepositoryName,
			)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil

}
