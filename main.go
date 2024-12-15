package main

import (
	"github.com/vitalii-komenda/got/controllers"
	"github.com/vitalii-komenda/got/postgres"
)

type AllControllers struct {
	CharactersController controllers.CharactersController
	SearchController     controllers.SearchController
}

func main() {
	db, err := postgres.NewDBPool()
	if err != nil {
		panic(err)
	}

	actorsRepo := postgres.NewActorsRepository(db)
	characterRepo := postgres.NewCharacterRepository(db, actorsRepo)
	relationshipsRepo := postgres.NewRelationshipsRepository(db, characterRepo)
	charactersController := controllers.NewCharactersController(characterRepo, relationshipsRepo)
	searchController := controllers.NewSearchController(characterRepo)

	allControllers := AllControllers{
		CharactersController: *charactersController,
		SearchController:     *searchController,
	}

	r := setupRouter(allControllers)
	r.Run(":8080")
}
