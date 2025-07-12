#!/bin/bash

# Usage: ./start-local-db.sh [container_name] [mongo_port] [data_volume]
CONTAINER_NAME="${1:-cosmos-mongo}"
MONGO_PORT="${2:-27017}"
DATA_VOLUME="${3:-cosmos-mongo-data}"

# We check if the container is already running
if [ "$(docker ps -q -f name="$CONTAINER_NAME")" ]; then
  echo "MongoDB container '$CONTAINER_NAME' is already running."
elif [ "$(docker ps -aq -f name="$CONTAINER_NAME")" ]; then
  echo "Starting existing MongoDB container '$CONTAINER_NAME'..."
  docker start "$CONTAINER_NAME"
  echo "MongoDB container '$CONTAINER_NAME' started on port $MONGO_PORT."
else
  # We create volume if it doesn't exist
  docker volume inspect "$DATA_VOLUME" > /dev/null 2>&1 || docker volume create "$DATA_VOLUME"
  # We run the MongoDB container
  docker run -d \
    --name "$CONTAINER_NAME" \
    -p "$MONGO_PORT":27017 \
    -v "$DATA_VOLUME":/data/db \
    mongo:latest
  echo "MongoDB container '$CONTAINER_NAME' started on port $MONGO_PORT."
fi
