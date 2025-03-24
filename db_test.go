package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		want task
	}{
		{
			want: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  "todo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.want.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.insert(tt.want.Name, tt.want.Project); err != nil {
				t.Fatalf("unable to insert tasks: %v", err)

			}

			tasks, err := tDB.getTasks()
			if err != nil {
				t.Fatalf("unable to get tasks: %v", err)
			}

			tt.want.Created = tasks[0].Created
			if !reflect.DeepEqual(tasks[0], tt.want) {
				t.Fatalf("want %v, got %v", tt.want, tasks[0])
			}
			if err := tDB.delete(tasks[0].ID); err != nil {
				t.Fatalf("unable to delete task: %v", err)
			}

			tasks, err = tDB.getTasks()
			if err != nil {
				t.Fatalf("unable to get tasks: %v", err)
			}
			if len(tasks) != 0 {
				t.Fatalf("expected tasks to be empty, got: %v", tasks)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	tests := []struct {
		want task
	}{
		{
			want: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.want.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.insert(tt.want.Name, tt.want.Project); err != nil {
				t.Fatalf("unable to insert tasks: %v", err)
			}

			task, err := tDB.getTask(1)
			if err != nil {
				t.Fatalf("unable to get task: %v", err)
			}

			tt.want.Created = task.Created
			if !reflect.DeepEqual(task, tt.want) {
				t.Fatalf("want %v, got %v", tt.want, task)
			}
		})
	}
}

func TestGetTasksByStatus(t *testing.T) {
	tests := []struct {
		want task
	}{
		{
			want: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.want.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.insert(tt.want.Name, tt.want.Project); err != nil {
				t.Fatalf("unable to insert tasks: %v", err)
			}

			tasks, err := tDB.getTasksByStatus(tt.want.Status)
			if err != nil {
				t.Fatalf("unable to get tasks: %v", err)
			}
			if len(tasks) < 1 {
				t.Fatalf("expected at least one task, got: %v", tasks)
			}

			tt.want.Created = tasks[0].Created
			if !reflect.DeepEqual(tasks[0], tt.want) {
				t.Fatalf("want %v, got %v", tt.want, tasks[0])
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		new  *task
		old  *task
		want task
	}{
		{
			new: &task{
				ID:      1,
				Name:    "strawberries",
				Project: "",
				Status:  "",
			},
			old: &task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
			want: task{
				ID:      1,
				Name:    "strawberries",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.new.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.insert(tt.old.Name, tt.old.Project); err != nil {
				t.Fatalf("unable to insert tasks: %v", err)
			}
			if err := tDB.update(*tt.new); err != nil {
				t.Fatalf("unable to update task: %v", err)
			}

			task, err := tDB.getTask(tt.want.ID)
			if err != nil {
				t.Fatalf("unable to get task: %v", err)
			}

			tt.want.Created = task.Created
			if !reflect.DeepEqual(task, tt.want) {
				t.Fatalf("want %v, got %v", tt.want, task)
			}
		})
	}
}

func setup() *taskDB {
	path := filepath.Join(os.TempDir(), "test.db")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	t := taskDB{
		db:      db,
		dataDir: path,
	}
	// if !t.tableExists("tasks") {
	err = t.createTable()
	if err != nil {
		log.Fatal(err)
	}
	// }

	return &t
}

func teardown(t *taskDB) {
	t.db.Close()
	os.Remove(t.dataDir)
}
