#!/bin/bash
set -x

GO=$(which go)
DOCKER=$(which docker)
GO_PROMPT_ENABLE_LOG="true"

# Clean up previous build
make clean

# Start the redis container
echo -n "Starting redis container...."
$DOCKER run --rm  --name some-redis -d -p 6379:6379 -v ~/redis-data:/data  redis redis-server --requirepass 'rpasswd' --appendonly yes
echo "Done"

# Build the executable
echo -n "Building redical...."
make build
echo "Done"

# Start redical
echo "Start redical"
export GO_PROMPT_ENABLE_LOG='true'
./redical --debug -P rpasswd

# Remove docker container
echo -n "Stopping redis container...."
$DOCKER stop some-redis
echo "Done"
