package internal

import (
	"database/sql"
	"log"
)

func InitDbSchema(db *sql.DB) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS todos(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT,
		completed INTEGER NOT NULL,
		createdAt INTEGER NOT NULL
	)`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table todos is ready to go!")
}
