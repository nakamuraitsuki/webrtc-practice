#!/bin/sh

# Run migrations
echo "Running migrations..."
cd ./backend
go run ./cmd/migration/main.go
cd ..
