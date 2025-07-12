#!/bin/bash

set -e

echo "Building..."
make all

echo "Running tests..."
make test

echo "Starting server..."
./service-catalog-api --db-driver=sqlite3 --db-dsn=services.db --port=8080
