#!/bin/sh

set -e

echo "run migrations"
source /app/app.env
/app/migrate -path /app/migrations -database "$DB_SOURCE" --verbose up

echo "start app"
exec "$@"
