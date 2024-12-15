package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vitalii-komenda/got/utils"
)

func NewDBPool() (*pgxpool.Pool, error) {
	pgHost := utils.MustGetEnvOrPanic("DB_HOST")
	pgUser := utils.MustGetEnvOrPanic("DB_USER")
	pgPass := utils.MustGetEnvOrPanic("DB_PASS")
	pgName := utils.MustGetEnvOrPanic("DB_NAME")
	pgPort, err := strconv.Atoi(utils.MustGetEnvOrPanic("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("error converting DB_PORT to int: %w", err)
	}

	dbUri := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", pgUser, pgPass, pgHost, pgPort, pgName)
	dbPool, err := pgxpool.Connect(context.Background(),
		dbUri,
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the DB: %w", err)
	}

	return dbPool, nil
}
