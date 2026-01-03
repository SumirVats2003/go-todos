package cmd

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type App struct {
	Db *sql.DB
}

func InitApp() App {
	db, err := sql.Open("sqlite3", "./data/todos.db")

	if err != nil {
		log.Fatal(err)
	}

	app := App{Db: db}
	defer db.Close()
	initDbSchema(db)

	return app
}

func initDbSchema(db *sql.DB) {
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
