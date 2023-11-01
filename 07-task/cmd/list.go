package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := ListTasks(false)

		if err != nil {
			exitf("%v", err)
		}

		if len(tasks) == 0 {
			fmt.Println("You have no tasks to complete! Why not take a vacation?")
			os.Exit(0) // no error happened
		}

		fmt.Println("You have the following tasks:")

		for _, task := range tasks {
			fmt.Printf("%d. %s\n", task.ID, task.Details)
		}

	},
}
