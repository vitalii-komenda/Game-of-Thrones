package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/vitalii-komenda/got/entities"
)

type CharsetTestSuite struct {
	suite.Suite
	repo   *CharactersRepository
	dbpool *pgxpool.Pool
	tx     *pgx.Tx
}

func (s *CharsetTestSuite) SetupSuite() {
	var err error

	s.dbpool, err = NewDBPool()
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *CharsetTestSuite) TearDownSuite() {
	s.dbpool.Close()
}

func (s *CharsetTestSuite) SetupTest() {
	ctx := context.Background()
	tx, err := s.dbpool.Begin(ctx)
	if err != nil {
		s.FailNow("Failed to begin transaction")
	}
	s.tx = &tx
	actorsRepo := NewActorsRepository(s.dbpool).WithTX(s.tx)
	s.repo = NewCharacterRepository(s.dbpool, actorsRepo).WithTX(s.tx)
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

func (s *CharsetTestSuite) TearDownTest() {
	if s.tx != nil {
		tx := *s.tx
		tx.Rollback(context.Background())
	}
}

func (s *CharsetTestSuite) TestDelete() {
	ctx := context.Background()

	var count int
	err := (*s.tx).QueryRow(ctx, "SELECT count(*) FROM characters").Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(1, count)
	err = s.repo.Delete(ctx, "Test Character")
	s.Require().NoError(err)

	err = (*s.tx).QueryRow(ctx, "SELECT count(*) FROM characters").Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(0, count)
}

func (s *CharsetTestSuite) TestGetCharacterID() {
	ctx := context.Background()
	characterName := "Test Character"

	id, err := s.repo.GetCharacterID(ctx, characterName)
	s.Require().NoError(err)
	s.Require().Equal(1, id)
}

func (s *CharsetTestSuite) TestCreateCharacter() {
	ctx := context.Background()
	characterEntryEntry := entities.CharacterEntry{
		CharacterName:       "Test Character 2",
		CharacterImageThumb: "http://test.com/thumb",
		CharacterImageFull:  "http://test.com/full",
		CharacterLink:       "http://test.com/character",
		Nickname:            "Test Nickname",
		Royal:               true,
	}
	houseNames := "Test House 1,Test House 2"

	id, err := s.repo.CreateCharacter(ctx, &characterEntryEntry, houseNames)
	s.Require().NoError(err)
	s.Require().Greater(id, 1)

	var nickname string
	err = (*s.tx).QueryRow(ctx, "SELECT nickname FROM characters WHERE character_name = 'Test Character 2'").Scan(&nickname)
	s.Require().NoError(err)
	s.Require().Equal("Test Nickname", nickname)
}

func (s *CharsetTestSuite) TestCreateCharacterAndActor() {
	ctx := context.Background()
	characterEntryEntry := entities.CharacterEntry{
		CharacterName:       "Test Character 3",
		CharacterImageThumb: "http://test.com/thumb",
		CharacterImageFull:  "http://test.com/full",
		CharacterLink:       "http://test.com/character",
		Nickname:            "Test Nickname",
		Royal:               true,
		ActorName:           "Test Actor",
		ActorLink:           "http://test.com/actor",
		HouseName:           []string{"Test House 1", "Test House 2"},
	}

	err := s.repo.CreateCharacterAndActor(ctx, &characterEntryEntry)
	s.Require().NoError(err)

	var characterId int
	err = (*s.tx).QueryRow(ctx, "SELECT character_id FROM characters WHERE character_name = 'Test Character 3'").Scan(&characterId)
	s.Require().NoError(err)
	s.Require().Greater(characterId, 1)

	var actorId int
	err = (*s.tx).QueryRow(ctx, "SELECT actor_id FROM actors WHERE actor_name = 'Test Actor'").Scan(&actorId)
	s.Require().NoError(err)
	s.Require().Greater(actorId, 0)

	var linkedActorId int
	err = (*s.tx).QueryRow(ctx, "SELECT actor_id FROM characters_actors WHERE character_id = $1", characterId).Scan(&linkedActorId)
	s.Require().NoError(err)
	s.Require().Equal(linkedActorId, actorId)
}

func (s *CharsetTestSuite) TestUpdateCharacterAndActor() {
	ctx := context.Background()
	characterEntryEntry := entities.CharacterEntry{
		CharacterName:       "Test Character",
		CharacterImageThumb: "http://test.com/thumb2",
		CharacterImageFull:  "http://test.com/full2",
		CharacterLink:       "http://test.com/character2",
		Nickname:            "Test Nickname2",
		Royal:               false,
		ActorName:           "Test Actor2",
		ActorLink:           "http://test.com/actor2",
		HouseName:           []string{"Test House 1", "Test House 2"},
	}

	id, err := s.repo.UpdateCharacterAndActor(ctx, &characterEntryEntry, "Test Character")
	s.Require().NoError(err)
	s.Require().Equal(1, id)

	var nickname string
	var royal bool
	err = (*s.tx).QueryRow(ctx, "SELECT nickname, royal FROM characters WHERE character_id = 1").Scan(&nickname, &royal)
	s.Require().NoError(err)
	s.Require().Equal("Test Nickname2", nickname)

	s.Require().NoError(err)
	s.Require().False(royal)
}

func TestRunCharsetTestSuite(t *testing.T) {
	suite.Run(t, &CharsetTestSuite{})
}
