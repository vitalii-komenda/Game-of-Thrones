package entities

import (
	"context"
	"encoding/json"
	"fmt"
)

type HouseNameType []string

type CharacterEntry struct {
	CharacterID         int           `json:"characterID,omitempty" db:"character_id"`
	CharacterName       string        `json:"characterName" db:"character_name"`
	HouseName           HouseNameType `json:"houseName,omitempty" db:"house_name"`
	CharacterImageThumb string        `json:"characterImageThumb,omitempty" db:"character_image_thumb"`
	CharacterImageFull  string        `json:"characterImageFull,omitempty" db:"character_image_full"`
	CharacterLink       string        `json:"characterLink,omitempty" db:"character_link"`
	ActorName           string        `json:"actorName,omitempty" db:"actor_name"`
	ActorLink           string        `json:"actorLink,omitempty" db:"actor_link"`
	Nickname            string        `json:"nickname,omitempty" db:"nickname"`
	Royal               bool          `json:"royal,omitempty" db:"royal"`
	Parents             []string      `json:"parents,omitempty" db:"parents"`
	Siblings            []string      `json:"siblings,omitempty" db:"siblings"`
	KilledBy            []string      `json:"killedBy,omitempty" db:"killed_by"`
	Killed              []string      `json:"killed,omitempty" db:"killed"`
	MarriedEngaged      []string      `json:"marriedEngaged,omitempty" db:"married_engaged"`
}

func (h *HouseNameType) UnmarshalJSON(data []byte) error {
	var singleName string
	if err := json.Unmarshal(data, &singleName); err == nil {
		*h = HouseNameType{singleName}
		return nil
	}

	var multipleNames []string
	if err := json.Unmarshal(data, &multipleNames); err == nil {
		*h = HouseNameType(multipleNames)
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into HouseNameType", string(data))
}

type CharacterEntryElastic struct {
	CharacterID         int           `json:"character_id,omitempty"`
	CharacterName       string        `json:"character_name"`
	HouseName           HouseNameType `json:"house_name,omitempty"`
	CharacterImageThumb string        `json:"character_image_thumb,omitempty"`
	CharacterImageFull  string        `json:"character_image_full,omitempty"`
	CharacterLink       string        `json:"character_link,omitempty"`
	ActorName           string        `json:"actor_name,omitempty"`
	ActorLink           string        `json:"actor_link,omitempty"`
	Nickname            string        `json:"nickname,omitempty"`
	Royal               bool          `json:"royal,omitempty"`
	Parents             []string      `json:"parents,omitempty"`
	Siblings            []string      `json:"siblings,omitempty"`
	KilledBy            []string      `json:"killed_by,omitempty"`
	Killed              []string      `json:"killed,omitempty"`
	MarriedEngaged      []string      `json:"married_engaged,omitempty"`
}

//go:generate moq -out ./../mocks/characters_repository.go -pkg mocks . CharactersRepository
type CharactersRepository interface {
	UpdateCharacterAndActor(ctx context.Context, characterEntryEntry *CharacterEntry, characterName string) (int, error)
	Delete(ctx context.Context, name string) error
	Get(ctx context.Context, name string) ([]CharacterEntry, error)
	GetAll(ctx context.Context, page int) ([]CharacterEntry, error)
	GetCharacterID(ctx context.Context, characterName string) (int, error)
	CreateCharacter(ctx context.Context, characterEntryEntry *CharacterEntry, houseNames string) (int, error)
	CreateCharacterAndActor(ctx context.Context, characterEntry *CharacterEntry) error
}
