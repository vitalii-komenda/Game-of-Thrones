package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	entities "github.com/vitalii-komenda/got/entities"
)

var ErrActorsRepoPersistenceFailure = fmt.Errorf("actors repo persistence failure")

var _ entities.ActorsRepository = &ActorsRepository{}

func NewActorsRepository(dbpool *pgxpool.Pool) *ActorsRepository {
	return &ActorsRepository{
		dbPool: dbpool,
	}
}

type ActorsRepository struct {
	dbPool *pgxpool.Pool
	tx     *pgx.Tx
}

func (r *ActorsRepository) WithTX(tx *pgx.Tx) *ActorsRepository {
	return &ActorsRepository{
		dbPool: r.dbPool,
		tx:     tx,
	}
}

func (r *ActorsRepository) getExecutor() PGXExecutor {
	if r.tx != nil {
		return *r.tx
	}
	return r.dbPool
}

func (r *ActorsRepository) Create(ctx context.Context, actorName string, actorLink string) (int, error) {
	sql, args, err := Psql.
		Insert("actors").
		Columns("actor_name", "actor_link").
		Values(actorName, actorLink).
		Suffix("RETURNING actor_id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrActorsRepoPersistenceFailure, err)
	}
	var id int
	err = r.getExecutor().QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrActorsRepoPersistenceFailure, err)
	}
	return id, nil
}

func (r *ActorsRepository) GetActorID(ctx context.Context, actorName string) (int, error) {
	sql, args, err := Psql.
		Select("actor_id").
		From("actors").
		Where("actor_name = ?", actorName).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("error building sql: %w", err)
	}
	row := r.getExecutor().QueryRow(ctx, sql, args...)
	var actorId int
	err = row.Scan(&actorId)
	if err != nil {
		return 0, fmt.Errorf("error scanning row: %w", err)
	}
	return actorId, nil
}

func (r *ActorsRepository) LinkActorToCharacter(ctx context.Context, actorId int, characterId int) error {
	sql, args, err := Psql.
		Insert("characters_actors").
		Columns("actor_id", "character_id").
		Values(actorId, characterId).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCharacterRepoPersistenceFailure, err)
	}
	_, err = r.getExecutor().Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCharacterRepoPersistenceFailure, err)
	}
	return nil
}

func (r *ActorsRepository) UnlinkActorFromCharacter(ctx context.Context, characterId int) error {
	sql, args, err := Psql.
		Delete("characters_actors").
		Where("character_id = ?", characterId).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCharacterRepoPersistenceFailure, err)
	}
	_, err = r.getExecutor().Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCharacterRepoPersistenceFailure, err)
	}
	return nil
}
