package db

import (
	"fmt"

	"github.com/e2b-dev/infra/packages/api/internal/api"
	"github.com/e2b-dev/infra/packages/api/internal/db/ent"
	"github.com/e2b-dev/infra/packages/api/internal/db/ent/env"
	"github.com/e2b-dev/infra/packages/api/internal/db/ent/envalias"

	"github.com/google/uuid"
)

func (db *DB) DeleteEnv(envID string) error {
	_, err := db.
		Client.
		Env.
		Delete().
		Where(env.ID(envID)).
		Exec(db.ctx)
	if err != nil {
		return fmt.Errorf("failed to delete env '%s': %w", envID, err)
	}

	return nil
}

func (db *DB) GetEnvs(teamID string) (result []*api.Environment, err error) {
	id, err := uuid.Parse(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse teamID: %w", err)
	}

	envs, err := db.
		Client.
		Env.
		Query().
		Where(env.Or(env.TeamID(id), env.Public(true))).
		WithEnvAliases().
		All(db.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list envs: %w", err)
	}

	for _, item := range envs {
		aliases := make([]string, len(item.Edges.EnvAliases))
		for i, item := range item.Edges.EnvAliases {
			aliases[i] = item.Alias
		}

		result = append(result, &api.Environment{
			EnvID:   item.ID,
			BuildID: item.BuildID.String(),
			Public:  item.Public,
			Aliases: &aliases,
		})
	}

	return result, nil
}

var ErrEnvNotFound = fmt.Errorf("env not found")

func (db *DB) GetEnv(aliasOrEnvID string, teamID string, canBePublic bool) (result *api.Environment, err error) {
	id, err := uuid.Parse(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse teamID: %w", err)
	}

	dbEnv, err := db.
		Client.
		Env.
		Query().
		Where(env.Or(env.HasEnvAliasesWith(envalias.Alias(aliasOrEnvID)), env.ID(aliasOrEnvID)), env.Or(env.TeamID(id), env.Public(true))).
		WithEnvAliases().
		Only(db.ctx)
	if err != nil {
		return nil, ErrEnvNotFound
	}

	if !canBePublic && dbEnv.TeamID != id {
		return nil, fmt.Errorf("you don't have access to this env '%s'", aliasOrEnvID)
	}

	aliases := make([]string, len(dbEnv.Edges.EnvAliases))
	for i, item := range dbEnv.Edges.EnvAliases {
		aliases[i] = item.Alias
	}

	return &api.Environment{
		EnvID:   dbEnv.ID,
		BuildID: dbEnv.BuildID.String(),
		Public:  dbEnv.Public,
		Aliases: &aliases,
	}, nil
}

func (db *DB) UpsertEnv(teamID, envID, buildID, dockerfile string) error {
	teamUUID, err := uuid.Parse(teamID)
	if err != nil {
		return fmt.Errorf("failed to parse teamID: %w", err)
	}

	buildUUID, err := uuid.Parse(buildID)
	if err != nil {
		return fmt.Errorf("failed to parse teamID: %w", err)
	}

	err = db.
		Client.
		Env.
		Create().
		SetID(envID).
		SetBuildID(buildUUID).
		SetTeamID(teamUUID).
		SetDockerfile(dockerfile).
		SetPublic(false).
		OnConflictColumns(env.FieldID).
		UpdateBuildID().
		UpdateDockerfile().
		UpdateUpdatedAt().
		Update(func(e *ent.EnvUpsert) {
			e.AddBuildCount(1)
		}).
		Exec(db.ctx)

	if err != nil {
		errMsg := fmt.Errorf("failed to upsert env with id '%s': %w", envID, err)

		fmt.Println(errMsg.Error())

		return errMsg
	}

	return nil
}

func (db *DB) HasEnvAccess(envID string, teamID string, canBePublic bool) (bool, error) {
	_, err := db.GetEnv(envID, teamID, canBePublic)
	if err != nil {
		return false, fmt.Errorf("failed to get env '%s': %w", envID, err)
	}

	return true, nil
}
