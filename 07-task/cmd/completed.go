package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completedCmd)
}

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List all of your complete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := ListTasks(true)

		if err != nil {
			exitf("%v", err)
		}

		if len(tasks) == 0 {
			fmt.Println("You have no completed tasks")
			os.Exit(0) // no error happened
		}

		fmt.Println("You have finished the following tasks today:")

		for _, task := range tasks {
			fmt.Printf("%d. %s\n", task.ID, task.Details)
		}

	},
}
