package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type ActorsTestSuite struct {
	suite.Suite
	repo   *ActorsRepository
	dbpool *pgxpool.Pool
	tx     *pgx.Tx
}

func (s *ActorsTestSuite) SetupSuite() {
	var err error

	s.dbpool, err = NewDBPool()
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ActorsTestSuite) TearDownSuite() {
	s.dbpool.Close()
}

func (s *ActorsTestSuite) SetupTest() {
	ctx := context.Background()
	tx, err := s.dbpool.Begin(ctx)
	if err != nil {
		s.FailNow("Failed to begin transaction")
	}
	s.tx = &tx
	s.repo = NewActorsRepository(s.dbpool).WithTX(s.tx)
	tx.Exec(ctx, `INSERT INTO characters
	 (
		character_id,
		character_name,
		house_name,
		character_image_thumb,
		character_image_full,
		character_link,
		nickname,
		royal
	) 
	 VALUES (
	 1,
	 'Test Character',
	 'Test House',
	 'http://test.com/thumb',
	 'http://test.com/full',
	 'http://test.com/character',
	 'Test Nickname',
	 true
	 )`)
}

func (s *ActorsTestSuite) TearDownTest() {
	if s.tx != nil {
		tx := *s.tx
		tx.Rollback(context.Background())
	}
}

func (s *ActorsTestSuite) TestCreateActor() {
	ctx := context.Background()
	actorName := "Test Actor"
	actorLink := "http://testactor.com"

	id, err := s.repo.Create(ctx, actorName, actorLink)
	s.Require().NoError(err)
	s.Require().Greater(id, 0)
}

func (s *ActorsTestSuite) TestGetActorID() {
	ctx := context.Background()
	actorName := "Test Actor"
	actorLink := "http://testactor.com"

	id, err := s.repo.Create(ctx, actorName, actorLink)
	s.Require().NoError(err)

	retrievedID, err := s.repo.GetActorID(ctx, actorName)
	s.Require().NoError(err)
	s.Require().Equal(id, retrievedID)
}

func (s *ActorsTestSuite) TestLinkActorToCharacter() {
	ctx := context.Background()
	actorName := "Test Actor"
	actorLink := "http://testactor.com"
	characterId := 1

	actorId, err := s.repo.Create(ctx, actorName, actorLink)
	s.Require().NoError(err)

	err = s.repo.LinkActorToCharacter(ctx, actorId, characterId)
	s.Require().NoError(err)

	var exists bool
	err = (*s.tx).QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM characters_actors WHERE actor_id=$1 AND character_id=$2)", actorId, characterId).Scan(&exists)
	s.Require().NoError(err)
	s.Require().True(exists)
}

func (s *ActorsTestSuite) TestUnlinkActorFromCharacter() {
	ctx := context.Background()
	actorName := "Test Actor"
	actorLink := "http://testactor.com"
	characterId := 1

	actorId, err := s.repo.Create(ctx, actorName, actorLink)
	s.Require().NoError(err)

	err = s.repo.LinkActorToCharacter(ctx, actorId, characterId)
	s.Require().NoError(err)

	err = s.repo.UnlinkActorFromCharacter(ctx, characterId)
	s.Require().NoError(err)

	var exists bool
	err = (*s.tx).QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM characters_actors WHERE actor_id=$1 AND character_id=$2)", actorId, characterId).Scan(&exists)
	s.Require().NoError(err)
	s.Require().False(exists)
}

func TestRunActorsTestSuite(t *testing.T) {
	suite.Run(t, &ActorsTestSuite{})
}
