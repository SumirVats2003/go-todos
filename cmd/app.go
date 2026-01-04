package cmd

import (
	"database/sql"
	"log"

	"github.com/SumirVats2003/go-todo/internal"
	_ "github.com/mattn/go-sqlite3"
)

func InitApp() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/todos.db")

	if err != nil {
		log.Fatal(err)
	}

	internal.InitDbSchema(db)

	repo := internal.InitRepository(db)
	InitTeaApp(repo)

	return db
}
