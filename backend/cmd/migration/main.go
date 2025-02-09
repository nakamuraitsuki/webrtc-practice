package  main

import (
	"log"
	"example.com/webrtc-practice/internal/infrastructure/repository/sqlite3"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlx.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	err = sqlite3.MigrateUser(db)
	if err != nil {
		log.Fatal(err)
	}
}