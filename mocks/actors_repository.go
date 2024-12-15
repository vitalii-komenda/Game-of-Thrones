// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/vitalii-komenda/got/entities"
	"sync"
)

// Ensure, that ActorsRepositoryMock does implement entities.ActorsRepository.
// If this is not the case, regenerate this file with moq.
var _ entities.ActorsRepository = &ActorsRepositoryMock{}

// ActorsRepositoryMock is a mock implementation of entities.ActorsRepository.
//
//	func TestSomethingThatUsesActorsRepository(t *testing.T) {
//
//		// make and configure a mocked entities.ActorsRepository
//		mockedActorsRepository := &ActorsRepositoryMock{
//			CreateFunc: func(ctx context.Context, actorName string, actorLink string) (int, error) {
//				panic("mock out the Create method")
//			},
//			GetActorIDFunc: func(ctx context.Context, actorName string) (int, error) {
//				panic("mock out the GetActorID method")
//			},
//			LinkActorToCharacterFunc: func(ctx context.Context, actorId int, characterId int) error {
//				panic("mock out the LinkActorToCharacter method")
//			},
//		}
//
//		// use mockedActorsRepository in code that requires entities.ActorsRepository
//		// and then make assertions.
//
//	}
type ActorsRepositoryMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(ctx context.Context, actorName string, actorLink string) (int, error)

	// GetActorIDFunc mocks the GetActorID method.
	GetActorIDFunc func(ctx context.Context, actorName string) (int, error)

	// LinkActorToCharacterFunc mocks the LinkActorToCharacter method.
	LinkActorToCharacterFunc func(ctx context.Context, actorId int, characterId int) error

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ActorName is the actorName argument value.
			ActorName string
			// ActorLink is the actorLink argument value.
			ActorLink string
		}
		// GetActorID holds details about calls to the GetActorID method.
		GetActorID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ActorName is the actorName argument value.
			ActorName string
		}
		// LinkActorToCharacter holds details about calls to the LinkActorToCharacter method.
		LinkActorToCharacter []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ActorId is the actorId argument value.
			ActorId int
			// CharacterId is the characterId argument value.
			CharacterId int
		}
	}
	lockCreate               sync.RWMutex
	lockGetActorID           sync.RWMutex
	lockLinkActorToCharacter sync.RWMutex
}

// Create calls CreateFunc.
func (mock *ActorsRepositoryMock) Create(ctx context.Context, actorName string, actorLink string) (int, error) {
	if mock.CreateFunc == nil {
		panic("ActorsRepositoryMock.CreateFunc: method is nil but ActorsRepository.Create was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		ActorName string
		ActorLink string
	}{
		Ctx:       ctx,
		ActorName: actorName,
		ActorLink: actorLink,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(ctx, actorName, actorLink)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedActorsRepository.CreateCalls())
func (mock *ActorsRepositoryMock) CreateCalls() []struct {
	Ctx       context.Context
	ActorName string
	ActorLink string
} {
	var calls []struct {
		Ctx       context.Context
		ActorName string
		ActorLink string
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// GetActorID calls GetActorIDFunc.
func (mock *ActorsRepositoryMock) GetActorID(ctx context.Context, actorName string) (int, error) {
	if mock.GetActorIDFunc == nil {
		panic("ActorsRepositoryMock.GetActorIDFunc: method is nil but ActorsRepository.GetActorID was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		ActorName string
	}{
		Ctx:       ctx,
		ActorName: actorName,
	}
	mock.lockGetActorID.Lock()
	mock.calls.GetActorID = append(mock.calls.GetActorID, callInfo)
	mock.lockGetActorID.Unlock()
	return mock.GetActorIDFunc(ctx, actorName)
}

// GetActorIDCalls gets all the calls that were made to GetActorID.
// Check the length with:
//
//	len(mockedActorsRepository.GetActorIDCalls())
func (mock *ActorsRepositoryMock) GetActorIDCalls() []struct {
	Ctx       context.Context
	ActorName string
} {
	var calls []struct {
		Ctx       context.Context
		ActorName string
	}
	mock.lockGetActorID.RLock()
	calls = mock.calls.GetActorID
	mock.lockGetActorID.RUnlock()
	return calls
}

// LinkActorToCharacter calls LinkActorToCharacterFunc.
func (mock *ActorsRepositoryMock) LinkActorToCharacter(ctx context.Context, actorId int, characterId int) error {
	if mock.LinkActorToCharacterFunc == nil {
		panic("ActorsRepositoryMock.LinkActorToCharacterFunc: method is nil but ActorsRepository.LinkActorToCharacter was just called")
	}
	callInfo := struct {
		Ctx         context.Context
		ActorId     int
		CharacterId int
	}{
		Ctx:         ctx,
		ActorId:     actorId,
		CharacterId: characterId,
	}
	mock.lockLinkActorToCharacter.Lock()
	mock.calls.LinkActorToCharacter = append(mock.calls.LinkActorToCharacter, callInfo)
	mock.lockLinkActorToCharacter.Unlock()
	return mock.LinkActorToCharacterFunc(ctx, actorId, characterId)
}

// LinkActorToCharacterCalls gets all the calls that were made to LinkActorToCharacter.
// Check the length with:
//
//	len(mockedActorsRepository.LinkActorToCharacterCalls())
func (mock *ActorsRepositoryMock) LinkActorToCharacterCalls() []struct {
	Ctx         context.Context
	ActorId     int
	CharacterId int
} {
	var calls []struct {
		Ctx         context.Context
		ActorId     int
		CharacterId int
	}
	mock.lockLinkActorToCharacter.RLock()
	calls = mock.calls.LinkActorToCharacter
	mock.lockLinkActorToCharacter.RUnlock()
	return calls
}