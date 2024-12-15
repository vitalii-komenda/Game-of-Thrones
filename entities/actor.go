package entities

import (
	"context"
)

type ActorEntry struct {
	ActorName string `json:"actorName,omitempty" db:"actor_name"`
	ActorLink string `json:"actorLink,omitempty" db:"actor_link"`
}

//go:generate moq -out ./../mocks/actors_repository.go -pkg mocks . ActorsRepository
type ActorsRepository interface {
	Create(ctx context.Context, actorName string, actorLink string) (int, error)
	GetActorID(ctx context.Context, actorName string) (int, error)
	LinkActorToCharacter(ctx context.Context, actorId int, characterId int) error
}
