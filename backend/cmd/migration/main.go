package  main

import (
	"log"
	"example.com/webrtc-practice/internal/infrastructure/repository/sqlite3"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlx.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	err = sqlite3.MigrateUser(db)
	if err != nil {
		log.Fatal(err)
	}
}