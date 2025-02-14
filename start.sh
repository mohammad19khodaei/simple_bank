#!/bin/sh

set -e

echo "run migrations"
/app/migrate -path /app/migrations -database "$DB_SOURCE" --verbose up

echo "start app"
exec "$@"
