package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/cobra"
)

/*
	Database schema
	Task table
	Task{
		ID field int auto inrement primary key
		Content field string field
		Marked field boolean field
	}
*/

type Task struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Marked  bool   `json:"marked"`
}

func AddTaskCommand(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add todo task to the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) > 0 {
				query := fmt.Sprintf("INSERT INTO task (`content`) VALUES ('%s')", args[0])
				insert, err := db.Query(query)
				if err != nil {
					panic(err.Error())
				}
				defer insert.Close()

			}
			return err
		},
	}
}

// ViewAllTasksCommand shows all tasks of the user
func ViewAllTasksCommand(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "view all tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := db.Query("SELECT * from task")
			if err != nil {
				return err
			}
			for results.Next() {
				var task Task
				err = results.Scan(&task.Id, &task.Content, &task.Marked)
				if err != nil {
					return err
				}
				fmt.Printf("%d. %s (%v)\n", task.Id, task.Content, task.Marked)
			}
			return nil
		},
	}
}

func RemovetaskCommand(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove task from task list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				key := args[0]
				_, err := db.Query(fmt.Sprintf("DELETE FROM task WHERE id = %s", key))
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// MarkCompleteTask amarks a task as complete
func MarkCompleteTask(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "mark",
		Short: "mark a task as done",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				key := args[0]
				_, err := db.Query(fmt.Sprintf("UPDATE task SET marked = true WHERE id = %s", key))
				if err != nil {
					return err
				}

			}
			return nil
		},
	}
}

func main() {

	cmd := &cobra.Command{
		Use:   "taskcli",
		Short: "manage user tasks",
	}

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/tasks")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	cmd.AddCommand(AddTaskCommand(db))
	cmd.AddCommand(ViewAllTasksCommand(db))
	cmd.AddCommand(RemovetaskCommand(db))
	cmd.AddCommand(MarkCompleteTask(db))
	cmd.Execute()
}
