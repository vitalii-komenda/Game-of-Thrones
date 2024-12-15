#!/bin/bash

# Check if GOOSE_DBSTRING is not set
if [ -z "$GOOSE_DBSTRING" ]; then
    # Compile GOOSE_DBSTRING from other environment variables
    export GOOSE_DBSTRING="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
    echo "Compiled GOOSE_DBSTRING for DB $DB_NAME at '$DB_HOST:$DB_PORT'"
fi

# Run Goose migrations
goose -dir /app/migrations up
