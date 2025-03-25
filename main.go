package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
)

func setupPath() string {
	scope := gap.NewScope(gap.User, "tasks")
	dirs, err := scope.DataDirs()
	if err != nil {
		log.Fatal(err)
	}

	var taskDir string
	if len(dirs) > 0 {
		taskDir = dirs[0]
	} else {
		fmt.Println("No task directory found, using home directory")
		taskDir, _ = os.UserHomeDir()
	}
	if err := initTaskDir(taskDir); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Using task directory: %s\n", taskDir)
	return taskDir
}

func initTaskDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}

func openDB(path string) (*taskDB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/tasks.db", path))
	if err != nil {
		fmt.Print("error opening database")
		return nil, err
	}

	t := taskDB{db: db, dataDir: path}
	if !t.tableExists("tasks") {
		err := t.createTable()
		if err != nil {
			return nil, err
		}
	}

	fmt.Println("Database opened successfully")

	return &t, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
