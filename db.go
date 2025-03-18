package main

import (
	"database/sql"
	"fmt"
	"time"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return [...]string{"todo", "in progress", "done"}[s]
}

type task struct {
	ID      int
	Name    string
	Project string
	Status  string
	Created time.Time
}

func (t task) FilterValue() string {
	return t.Name
}

func (t task) Title() string {
	return t.Name
}

func (t task) Description() string {
	return t.Project
}

type taskDB struct {
	db      *sql.DB
	dataDir string
}

func (t *taskDB) tableExists(name string) bool {
	var query string = fmt.Sprintf("SELECT * FROM %s", name)

	if _, err := t.db.Query(query); err != nil {
		return true
	}
	return false
}

func (t *taskDB) createTable() error {
	_, err := t.db.Exec(`
	CREATE TABLE "tasks" ( 
		"id" INTEGER,
		"name" TEXT NOT NULL, 
		"project" TEXT, 
		"status" TEXT, 
		"created" DATETIME, 
		PRIMARY KEY("id" AUTOINCREMENT)
	)`)
	return err
}

func (t *taskDB) insert(name, project string) error {
	_, err := t.db.Exec(`
	INSERT INTO tasks (
		name, project, status, created
	) VALUES ( ?, ?, ?, ? )`,
		name,
		project,
		todo.String(),
		time.Now())

	return err
}

func (t *taskDB) update(task task) error {
	return nil
}

func (t *taskDB) delete(id uint) error {
	_, err := t.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}
