psql -U $POSTGRES_USER -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | \
    grep -q 1 | \
        psql -U $POSTGRES_USER -c "CREATE DATABASE $DB_NAME"


psql -U $POSTGRES_USER -tc "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}'" | \
    grep -q 1 | \
        psql -U $POSTGRES_USER -c "CREATE DATABASE ${DB_NAME}_test"
