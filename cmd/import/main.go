package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/vitalii-komenda/got/entities"
	"github.com/vitalii-komenda/got/postgres"
)

type Characters struct {
	Characters []entities.CharacterEntry `json:"characters"`
}

// DB_HOST=localhost DB_USER=postgres DB_PASS=postgres DB_NAME=got DB_PORT=5433 go run cmd/import/main.go
func main() {
	var err error
	ctx := context.Background()
	db, err := postgres.NewDBPool()
	if err != nil {
		panic(err)
	}

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		panic(err)
	}
	defer tx.Rollback(ctx)

	// characterDetails := postgres.NewCharacterDetailsRepository(db).WithTX(&tx)
	actorsRepo := postgres.NewActorsRepository(db).WithTX(&tx)
	characterRepo := postgres.NewCharacterRepository(db, actorsRepo).WithTX(&tx)
	relationshipsRepo := postgres.NewRelationshipsRepository(db, characterRepo).WithTX(&tx)

	file, err := os.Open("data/got-characters.json")
	if err != nil {
		log.Fatalf("Unable to open file: %v\n", err)
		panic(err)
	}
	defer file.Close()
	var characters Characters
	if err := json.NewDecoder(file).Decode(&characters); err != nil {
		log.Fatalf("Unable to decode JSON: %v\n", err)
		panic(err)
	}

	fmt.Print("Importing data...\n")

	// create characters and actors. relate them together
	for _, charData := range characters.Characters {
		fmt.Print("Creating character: ", charData.CharacterName, "\n")

		err = characterRepo.CreateCharacterAndActor(ctx, &charData)
		if err != nil {
			log.Fatalf("Unable to create character: %v-%v\n", err, charData.CharacterName)
			panic(err)
		}
	}

	// add the rest of the relationships
	for _, charData := range characters.Characters {
		fmt.Print("Adding relationships for: ", charData.CharacterName, "\n")

		err = relationshipsRepo.AddAll(ctx, charData)
		if err != nil {
			log.Fatalf("Unable to add relationships: %v\n", err)
			panic(err)
		}
	}

	// if err == nil {
	// 	err = characterDetails.Refresh(ctx)
	// 	if err != nil {
	// 		log.Fatalf("Unable to refresh materialized view: %v\n", err)
	// 		panic(err)
	// 	}
	// }

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("Unable to commit transaction: %v\n", err)
	}

	fmt.Print("\n\nImport completed\n\n")
}
