set -x
set -eo pipefail

if ! [ -x "$(command -v psql)" ]; then
    echo >&2 "Error: psql is not installed."
    exit 1
fi

if ! [ -x "$(command -v goose)" ]; then
    echo >&2 "Error: goose is not installed."
    echo >&2 "Use:"
    echo >&2 "go install github.com/pressly/goose/v3/cmd/goose@latest"
    echo >&2 "to install it."
    exit 1
fi

DB_HOST=${DATABASE_HOST:=localhost}
DB_USER=${DATABASE_USER:=postgres}
DB_PASSWORD="${DATABASE_PASSWORD:=password}"
DB_NAME="${DATABASE_NAME:=lingostruct}"
DB_PORT="${DATABASE_PORT:=5432}"

if [[ -z "${SKIP_DOCKER}" ]]
then
    docker run \
    -e POSTGRES_USER=${DB_USER} \
    -e POSTGRES_PASSWORD=${DB_PASSWORD} \
    -e POSTGRES_DB=${DB_NAME} \
    -p "${DB_PORT}":5432 \
    -d postgres \
    postgres -N 1000
fi

export PGPASSWORD="${DB_PASSWORD}"
until psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d "postgres" -c '\q'; do
    >&2 echo "Postgres is still unavailable - sleeping"
    sleep 1
done

>&2 echo "Postgres is up and running on ${DB_HOST}:${DB_PORT}!"

export GOOSE_MIGRATION_DIR=migrations
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}

export DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}
goose up

>&2 echo "Postgres has been migrated, ready to go!"
