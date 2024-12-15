package entities

import (
	"context"
)

//go:generate moq -out ./../mocks/relationships_repository.go -pkg mocks . RelationshipsRepository
type RelationshipsRepository interface {
	UpdateAll(ctx context.Context, character CharacterEntry) error
	AddAll(ctx context.Context, character CharacterEntry) error
	AddParent(ctx context.Context, characterID int, characterParentId int) error
	AddSibling(ctx context.Context, characterID int, characterSiblingId int) error
	AddKilled(ctx context.Context, characterID int, characterKilledId int) error
	AddMarriedEngaged(ctx context.Context, characterID int, characterMarriedEngagedId int) error
}
