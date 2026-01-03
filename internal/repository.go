package internal

import (
	"database/sql"
	"log"
	"time"

	"github.com/SumirVats2003/go-todo/internal/model"
)

type Repository struct {
	Store *sql.DB
}

func InitRepository(db *sql.DB) Repository {
	repo := Repository{Store: db}
	return repo
}

func (r Repository) GetAllTodos() []model.Todo {
	todos := make([]model.Todo, 0)
	sql := "SELECT * FROM todos"

	rows, err := r.Store.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo model.Todo
		err = rows.Scan(&todo.Id, &todo.Title, &todo.Content, &todo.Completed, &todo.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}

	return todos
}

func (r Repository) GetTodo(id int) model.Todo {
	var todo model.Todo

	sql := "SELECT * FROM todos WHERE id = ?"
	row := r.Store.QueryRow(sql, id)

	err := row.Scan(&todo.Id, &todo.Title, &todo.Content, &todo.Completed, &todo.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}

	return todo
}

func (r Repository) CreateTodo(todo model.Todo) {
	sql := `INSERT INTO todos 
	(id, title, content, completed, createdAt) 
	VALUES (?, ?, ?, ?, ?)`

	_, err := r.Store.Exec(sql, todo.Id, todo.Title, todo.Content, todo.Completed, time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}
}

func (r Repository) UpdateTodo(id int, updatedTodo model.Todo) error {
	sql := `UPDATE todos
	SET title = ?, content = ?, completed = ?
	WHERE id = ?`

	_, err := r.Store.Exec(sql, updatedTodo.Title, updatedTodo.Content, updatedTodo.Completed, id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r Repository) DeleteTodo(id int) error {
	sql := `DELETE FROM todos WHERE id = ?`

	_, err := r.Store.Exec(sql, id)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
