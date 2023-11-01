package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		// asci to int

		if len(args) == 0 {
			exitf("Please provide a task number to be removed")
		}

		taskID, err := strconv.Atoi(args[0])

		if err != nil {
			exitf("%v", err)
		}
		task := &Task{ID: taskID}
		if err := DeleteTask(task); err != nil {
			exitf("%v", err)
		}

		fmt.Printf(`You have deleted %q the task.`, task.Details)
	},
}
