package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
	Run: func(cmd *cobra.Command, args []string) {
		// asci to int

		if len(args) == 0 {
			exitf("Please provide a task number to mark as complete")
		}

		taskID, err := strconv.Atoi(args[0])

		if err != nil {
			exitf("%v", err)
		}
		task := &Task{ID: taskID}
		if err := MarkTaskAsCompleted(task); err != nil {
			exitf("%v", err)
		}

		fmt.Printf(`You have completed %q the task.`, task.Details)
	},
}
