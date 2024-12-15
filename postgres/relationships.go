package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vitalii-komenda/got/entities"
)

var ErrRelationshipsRepository = fmt.Errorf("relationships repository failure")

var _ entities.RelationshipsRepository = &RelationshipsRepository{}

func NewRelationshipsRepository(dbpool *pgxpool.Pool, charactersRepo *CharactersRepository) *RelationshipsRepository {
	return &RelationshipsRepository{
		dbPool:         dbpool,
		charactersRepo: charactersRepo,
	}
}

type RelationshipsRepository struct {
	dbPool         *pgxpool.Pool
	tx             *pgx.Tx
	charactersRepo *CharactersRepository
}

func (r *RelationshipsRepository) WithTX(tx *pgx.Tx) *RelationshipsRepository {
	return &RelationshipsRepository{
		dbPool:         r.dbPool,
		tx:             tx,
		charactersRepo: r.charactersRepo,
	}
}

func (r *RelationshipsRepository) getExecutor() PGXExecutor {
	if r.tx != nil {
		return *r.tx
	}
	return r.dbPool
}

func (r *RelationshipsRepository) AddRelationship(ctx context.Context, characterID int, characterRelationshipId int, relationshipType string) error {
	sql, args, err := Psql.
		Insert("relationships").
		Columns("character_id", "character_relationship_id", "relationship_type").
		Values(characterID, characterRelationshipId, relationshipType).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRelationshipsRepository, err)
	}
	_, err = r.getExecutor().Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRelationshipsRepository, err)
	}
	return nil
}

func (r *RelationshipsRepository) AddParent(ctx context.Context, characterID int, characterParentId int) error {
	return r.AddRelationship(ctx, characterID, characterParentId, "parent")
}

func (r *RelationshipsRepository) AddSibling(ctx context.Context, characterID int, characterSiblingId int) error {
	return r.AddRelationship(ctx, characterID, characterSiblingId, "sibling")
}

func (r *RelationshipsRepository) AddKilled(ctx context.Context, characterID int, characterKilledId int) error {
	return r.AddRelationship(ctx, characterID, characterKilledId, "killed")
}

func (r *RelationshipsRepository) AddMarriedEngaged(ctx context.Context, characterID int, characterMarriedEngagedId int) error {
	return r.AddRelationship(ctx, characterID, characterMarriedEngagedId, "married_engaged")
}

func (r *RelationshipsRepository) AddAll(ctx context.Context, character entities.CharacterEntry) error {
	characterId, err := r.charactersRepo.GetCharacterID(ctx, character.CharacterName)
	if err != nil || characterId == 0 {
		return fmt.Errorf("unable to get character %s %w: %v", character.CharacterName, ErrRelationshipsRepository, err)
	}

	// add siblings
	for _, sibling := range character.Siblings {
		siblingFromDB, err := r.charactersRepo.GetCharacterID(ctx, sibling)
		if err != nil || siblingFromDB == 0 {
			fmt.Printf("Sibling not found: %v\n", sibling)
			continue
		}

		err = r.AddSibling(ctx, characterId, siblingFromDB)
		if err != nil {
			return fmt.Errorf("unable to create sibling %s %w: %v", character.CharacterName, ErrRelationshipsRepository, err)
		}
	}

	// add parents
	for _, parent := range character.Parents {
		parentFromDB, err := r.charactersRepo.GetCharacterID(ctx, parent)
		if err != nil {
			return fmt.Errorf("unable to get parent %s %w: %v", parent, ErrRelationshipsRepository, err)
		}
		if parentFromDB == 0 {
			fmt.Printf("Parent not found: %v\n", parent)
			continue
		}

		err = r.AddParent(ctx, characterId, parentFromDB)
		if err != nil {
			return fmt.Errorf("unable to create parent %s %w: %v", character.CharacterName, ErrRelationshipsRepository, err)
		}
	}

	// add killed
	for _, killed := range character.Killed {
		killedFromDB, err := r.charactersRepo.GetCharacterID(ctx, killed)
		if err != nil || killedFromDB == 0 {
			fmt.Printf("Killed not found: %v\n", killed)
			continue
		}

		err = r.AddKilled(ctx, characterId, killedFromDB)
		if err != nil {
			return fmt.Errorf("unable to create killed %s %w: %v", character.CharacterName, ErrRelationshipsRepository, err)
		}
	}

	// add married_engaged
	for _, marriedEngaged := range character.MarriedEngaged {
		marriedEngagedFromDB, err := r.charactersRepo.GetCharacterID(ctx, marriedEngaged)
		if err != nil {
			return fmt.Errorf("unable to get marriedEngaged %s %w: %v", marriedEngaged, ErrRelationshipsRepository, err)
		}
		if marriedEngagedFromDB == 0 {
			fmt.Printf("marriedEngaged not found: %v\n", marriedEngaged)
			continue
		}

		err = r.AddMarriedEngaged(ctx, characterId, marriedEngagedFromDB)
		if err != nil {
			return fmt.Errorf("unable to create marriedEngaged %s %w: %v", character.CharacterName, ErrRelationshipsRepository, err)
		}
	}

	return nil
}

func (r *RelationshipsRepository) UpdateAll(ctx context.Context, character entities.CharacterEntry) error {
	characterId, err := r.charactersRepo.GetCharacterID(ctx, character.CharacterName)
	if err != nil || characterId == 0 {
		return fmt.Errorf("unable to get character %s %w: %v", character.CharacterName, ErrRelationshipsRepository, err)
	}

	// delete all relationships
	err = r.DeleteAll(ctx, characterId)
	if err != nil {
		return fmt.Errorf("unable to delete all relationships %w: %v", ErrRelationshipsRepository, err)
	}

	// add all relationships
	err = r.AddAll(ctx, character)
	if err != nil {
		return fmt.Errorf("unable to add all relationships %w: %v", ErrRelationshipsRepository, err)
	}

	return nil
}

func (r *RelationshipsRepository) DeleteAll(ctx context.Context, characterID int) error {
	sql, args, err := Psql.
		Delete("relationships").
		Where("character_id = ?", characterID).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRelationshipsRepository, err)
	}
	_, err = r.getExecutor().Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRelationshipsRepository, err)
	}

	return nil
}
