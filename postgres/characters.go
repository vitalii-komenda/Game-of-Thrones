package postgres

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	entities "github.com/vitalii-komenda/got/entities"
)

var ErrCharacterRepoPersistenceFailure = fmt.Errorf("character repo persistence failure")

var _ entities.CharactersRepository = &CharactersRepository{}

func NewCharacterRepository(dbpool *pgxpool.Pool, actorsRepo *ActorsRepository) *CharactersRepository {
	return &CharactersRepository{
		dbPool:     dbpool,
		actorsRepo: actorsRepo,
	}
}

type CharactersRepository struct {
	dbPool     *pgxpool.Pool
	tx         *pgx.Tx
	actorsRepo *ActorsRepository
}

func (r *CharactersRepository) WithTX(tx *pgx.Tx) *CharactersRepository {
	return &CharactersRepository{
		dbPool:     r.dbPool,
		tx:         tx,
		actorsRepo: r.actorsRepo,
	}
}

func (r *CharactersRepository) getExecutor() PGXExecutor {
	if r.tx != nil {
		return *r.tx
	}
	return r.dbPool
}

// Gets or creates an actor first, then creates character
func (r *CharactersRepository) CreateCharacterAndActor(ctx context.Context, characterEntryEntry *entities.CharacterEntry) error {
	var actorId int
	var err error
	var characterId int
	if characterEntryEntry.ActorName != "" {
		actorId, err = r.actorsRepo.GetActorID(ctx, characterEntryEntry.ActorName)
		if err != nil {
			actorId, err = r.actorsRepo.Create(ctx, characterEntryEntry.ActorName, characterEntryEntry.ActorLink)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrCharacterRepoPersistenceFailure, err)
			}
		}
	}

	existingCharacter, err := r.GetCharacterID(ctx, characterEntryEntry.CharacterName)
	if err != nil || existingCharacter == 0 {
		houseNames := strings.Join(characterEntryEntry.HouseName, ",")
		characterId, err = r.CreateCharacter(ctx, characterEntryEntry, houseNames)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrCharacterRepoPersistenceFailure, err)
		}
	} else {
		characterId = existingCharacter
	}

	if characterEntryEntry.ActorName != "" {
		return r.actorsRepo.LinkActorToCharacter(ctx, actorId, characterId)
	}
	return nil
}

func (r *CharactersRepository) UpdateCharacterAndActor(ctx context.Context, characterEntryEntry *entities.CharacterEntry, characterName string) (int, error) {
	existingCharacter, err := r.GetCharacterID(ctx, characterName)
	if err != nil {
		return 0, fmt.Errorf("GetCharacterID %w: %v", ErrCharacterRepoPersistenceFailure, err)
	}
	if existingCharacter == 0 {
		return 0, fmt.Errorf("%w: character not found", ErrCharacterRepoPersistenceFailure)
	}

	houseNames := strings.Join(characterEntryEntry.HouseName, ",")
	sql, args, err := Psql.
		Update("characters").
		Set("house_name", houseNames).
		Set("character_image_thumb", characterEntryEntry.CharacterImageThumb).
		Set("character_image_full", characterEntryEntry.CharacterImageFull).
		Set("character_link", characterEntryEntry.CharacterLink).
		Set("nickname", characterEntryEntry.Nickname).
		Set("royal", characterEntryEntry.Royal).
		Where("character_name = ?", characterName).
		Suffix("RETURNING character_id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrCharacterRepoPersistenceFailure, err)
	}
	var id int
	err = r.getExecutor().QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("update query %w: %v", ErrCharacterRepoPersistenceFailure, err)
	}

	// unlink actor from character
	err = r.actorsRepo.UnlinkActorFromCharacter(ctx, existingCharacter)
	if err != nil {
		return 0, fmt.Errorf("unlink actor %w: %v", ErrCharacterRepoPersistenceFailure, err)
	}

	var actorId int
	if characterEntryEntry.ActorName != "" {
		actorId, err = r.actorsRepo.GetActorID(ctx, characterEntryEntry.ActorName)
		if err != nil {
			actorId, err = r.actorsRepo.Create(ctx, characterEntryEntry.ActorName, characterEntryEntry.ActorLink)
			if err != nil {
				return 0, fmt.Errorf("create actor %w: %v", ErrCharacterRepoPersistenceFailure, err)
			}
		}

		err = r.actorsRepo.LinkActorToCharacter(ctx, actorId, existingCharacter)
		if err != nil {
			return 0, fmt.Errorf("link to character %w: %v", ErrCharacterRepoPersistenceFailure, err)
		}
	}

	return id, nil
}

func (r *CharactersRepository) CreateCharacter(ctx context.Context, characterEntryEntry *entities.CharacterEntry, houseNames string) (int, error) {
	sql, args, err := Psql.
		Insert("characters").
		Columns(
			"character_name",
			"house_name",
			"character_image_thumb",
			"character_image_full",
			"character_link",
			"nickname",
			"royal",
		).
		Values(
			characterEntryEntry.CharacterName,
			houseNames,
			characterEntryEntry.CharacterImageThumb,
			characterEntryEntry.CharacterImageFull,
			characterEntryEntry.CharacterLink,
			characterEntryEntry.Nickname,
			characterEntryEntry.Royal,
		).
		Suffix("RETURNING character_id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrCharacterRepoPersistenceFailure, err)
	}

	var id int
	err = r.getExecutor().QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrCharacterRepoPersistenceFailure, err)
	}
	return id, nil
}

func (r *CharactersRepository) GetCharacterID(ctx context.Context, characterName string) (int, error) {
	decodedName, err := url.QueryUnescape(characterName)
	if err != nil {
		return 0, fmt.Errorf("error decoding character name: %w", err)
	}

	sql, args, err := Psql.
		Select("character_id").
		From("characters").
		Where("character_name = ?", decodedName).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("error building sql: %w", err)
	}
	row := r.getExecutor().QueryRow(ctx, sql, args...)
	var characterId int
	err = row.Scan(&characterId)
	if err != nil {
		return 0, fmt.Errorf("error scanning row: %w", err)
	}
	return characterId, nil
}

func (r *CharactersRepository) Delete(ctx context.Context, name string) error {
	sql, args, err := Psql.
		Delete("characters").
		Where("character_name = ?", name).
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

func (r *CharactersRepository) Get(ctx context.Context, name string) ([]entities.CharacterEntry, error) {
	decodedName, err := url.QueryUnescape(name)
	if err != nil {
		return nil, fmt.Errorf("error decoding character name: %w", err)
	}
	var houseName pgtype.TextArray

	sql, args, err := Psql.
		Select(`
			c.character_id,
			c.character_name,
			string_to_array(c.house_name, '-') AS house_name,
			c.character_image_thumb,
			c.character_image_full,
			c.character_link,
			c.nickname,
			c.royal,
			COALESCE(a.actor_name, '') AS actor_name,
			COALESCE(a.actor_link, '') AS actor_link,
			COALESCE(array_remove(array_agg(DISTINCT CASE WHEN r.relationship_type = 'parent' THEN related_character.character_name END), NULL), '{}') AS parents,
			COALESCE(array_remove(array_agg(DISTINCT CASE WHEN r.relationship_type = 'sibling' THEN related_character.character_name END), NULL), '{}') AS siblings,
			COALESCE(array_remove(array_agg(DISTINCT CASE WHEN r.relationship_type = 'killed' THEN related_character.character_name END), NULL), '{}') AS killed,
			COALESCE(array_remove(array_agg(DISTINCT CASE WHEN r.relationship_type = 'killed_by' THEN related_character.character_name END), NULL), '{}') AS killed_by,
			COALESCE(array_remove(array_agg(DISTINCT CASE WHEN r.relationship_type = 'married_engaged' THEN related_character.character_name END), NULL), '{}') AS married_engaged
	`).
		From("characters AS c").
		LeftJoin("characters_actors AS ca ON c.character_id = ca.character_id").
		LeftJoin("actors AS a ON ca.actor_id = a.actor_id").
		LeftJoin("relationships AS r ON c.character_id = r.character_id").
		LeftJoin("characters AS related_character ON r.character_relationship_id = related_character.character_id").
		Where("c.character_name = ?", decodedName).
		GroupBy("c.character_id, a.actor_name, a.actor_link").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building sql: %w", err)
	}
	rows, err := r.getExecutor().Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()
	var characters []entities.CharacterEntry
	for rows.Next() {
		var c entities.CharacterEntry
		err := rows.Scan(
			&c.CharacterID,
			&c.CharacterName,
			&houseName,
			&c.CharacterImageThumb,
			&c.CharacterImageFull,
			&c.CharacterLink,
			&c.Nickname,
			&c.Royal,
			&c.ActorName,
			&c.ActorLink,
			pq.Array(&c.Parents),
			pq.Array(&c.Siblings),
			pq.Array(&c.Killed),
			pq.Array(&c.KilledBy),
			pq.Array(&c.MarriedEngaged),
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if houseName.Status == pgtype.Present {
			c.HouseName = make(entities.HouseNameType, len(houseName.Elements))
			for i, elem := range houseName.Elements {
				c.HouseName[i] = elem.String
			}
		} else {
			c.HouseName = entities.HouseNameType{}
		}

		characters = append(characters, c)
	}

	return characters, nil
}

func (r *CharactersRepository) GetAll(ctx context.Context, page int) ([]entities.CharacterEntry, error) {
	var houseName pgtype.TextArray

	sql, args, err := Psql.
		Select(`
			c.character_id,
			c.character_name,
			string_to_array(c.house_name, '-') AS house_name,
			c.character_image_thumb,
			c.character_image_full,
			c.character_link,
			c.nickname,
			c.royal,
			COALESCE(a.actor_name, '') AS actor_name,
			COALESCE(a.actor_link, '') AS actor_link,
			COALESCE(array_remove(array_agg(DISTINCT CASE WHEN r.relationship_type = 'parent' THEN related_character.character_name END), NULL), '{}') AS parents,
			COALESCE(array_remove(array_agg(DISTINCT CASE WHEN r.relationship_type = 'sibling' THEN related_character.character_name END), NULL), '{}') AS siblings,
			COALESCE(array_remove(array_agg(DISTINCT CASE WHEN r.relationship_type = 'killed' THEN related_character.character_name END), NULL), '{}') AS killed,
			COALESCE(array_remove(array_agg(DISTINCT CASE WHEN r.relationship_type = 'killed_by' THEN related_character.character_name END), NULL), '{}') AS killed_by,
			COALESCE(array_remove(array_agg(DISTINCT CASE WHEN r.relationship_type = 'married_engaged' THEN related_character.character_name END), NULL), '{}') AS married_engaged
	`).
		From("characters AS c").
		LeftJoin("characters_actors AS ca ON c.character_id = ca.character_id").
		LeftJoin("actors AS a ON ca.actor_id = a.actor_id").
		LeftJoin("relationships AS r ON c.character_id = r.character_id").
		LeftJoin("characters AS related_character ON r.character_relationship_id = related_character.character_id").
		GroupBy("c.character_id, a.actor_name, a.actor_link").
		Limit(25).
		Offset(uint64(page) * 25).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building sql: %w", err)
	}

	rows, err := r.getExecutor().Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()
	var characters []entities.CharacterEntry
	for rows.Next() {
		var c entities.CharacterEntry
		err := rows.Scan(
			&c.CharacterID,
			&c.CharacterName,
			&houseName,
			&c.CharacterImageThumb,
			&c.CharacterImageFull,
			&c.CharacterLink,
			&c.Nickname,
			&c.Royal,
			&c.ActorName,
			&c.ActorLink,
			pq.Array(&c.Parents),
			pq.Array(&c.Siblings),
			pq.Array(&c.Killed),
			pq.Array(&c.KilledBy),
			pq.Array(&c.MarriedEngaged),
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if houseName.Status == pgtype.Present {
			c.HouseName = make(entities.HouseNameType, len(houseName.Elements))
			for i, elem := range houseName.Elements {
				c.HouseName[i] = elem.String
			}
		} else {
			c.HouseName = entities.HouseNameType{}
		}

		characters = append(characters, c)
	}

	return characters, nil
}
