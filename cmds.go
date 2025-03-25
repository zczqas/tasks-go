package main

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "A CLI task manager",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var addCmd = &cobra.Command{
	Use:   "add NAME",
	Short: "Add a new task with optional project name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := openDB(setupPath())
		if err != nil {
			return err
		}

		defer t.db.Close()
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}
		if err := t.insert(args[0], project); err != nil {
			return err
		}
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete a task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := openDB(setupPath())
		if err != nil {
			return err
		}

		defer t.db.Close()
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		return t.delete(uint(id))
	},
}

var updateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := openDB(setupPath())
		if err != nil {
			return nil
		}

		defer t.db.Close()
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}

		prog, err := cmd.Flags().GetInt("status")
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		var status string
		switch prog {
		case int(inProgress):
			status = inProgress.String()
		case int(done):
			status = done.String()
		default:
			status = todo.String()
		}

		newTask := task{
			ID:      uint(id),
			Name:    name,
			Project: project,
			Status:  status,
			Created: time.Time{},
		}
		return t.update(newTask)
	},
}
