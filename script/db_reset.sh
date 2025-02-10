#!/bin/sh

# Reset database
echo "Resetting database..."
cd ./backend
rm database.db

# Run migrations
go run ./cmd/migration/main.go
cd ..