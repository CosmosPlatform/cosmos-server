#!/bin/bash

# Usage: ./start-local-db.sh [container_name] [pg_port] [data_volume]
CONTAINER_NAME="${1:-cosmos-postgres}"
PG_PORT="${2:-5432}"
DATA_VOLUME="${3:-cosmos-postgres-data}"
PG_PASSWORD="postgres"

# We check if the container is already running
if [ "$(docker ps -q -f name="$CONTAINER_NAME")" ]; then
  echo "PostgresSQL container '$CONTAINER_NAME' is already running."
elif [ "$(docker ps -aq -f name="$CONTAINER_NAME")" ]; then
  echo "Starting existing PostgresSQL container '$CONTAINER_NAME'..."
  docker start "$CONTAINER_NAME"
  echo "PostgresSQL container '$CONTAINER_NAME' started on port $PG_PORT."
else
  # We create volume if it doesn't exist
  docker volume inspect "$DATA_VOLUME" > /dev/null 2>&1 || docker volume create "$DATA_VOLUME"
  # We run the PostgresSQL container
  docker run -d \
    --name "$CONTAINER_NAME" \
    -p "$PG_PORT":5432 \
    -v "$DATA_VOLUME":/var/lib/postgresql/data \
    -e POSTGRES_PASSWORD="$PG_PASSWORD" \
    -e POSTGRES_DB="cosmos" \
    postgres:latest
  echo "PostgresSQL container '$CONTAINER_NAME' started on port $PG_PORT."
fi